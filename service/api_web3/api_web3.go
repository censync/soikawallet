package api_web3

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	addr = "127.0.0.1:8114"

	protocolVersion = 1
)

type Web3Connection struct {
	walletId string
	uiEvents *event_bus.EventBus
	w3Events *event_bus.EventBus
	server   *http.Server

	upgrader websocket.Upgrader // use default options
	done     chan bool
	wg       *sync.WaitGroup
	hub      map[string]*websocket.Conn
	accepted map[string]bool
	rejected map[string]bool
}

func NewWeb3Connection(cfg *config.Config, wg *sync.WaitGroup, uiEvents, w3Events *event_bus.EventBus) *Web3Connection {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	return &Web3Connection{
		walletId: ``,
		uiEvents: uiEvents,
		w3Events: w3Events,
		server:   server,
		upgrader: upgrader,
		done:     make(chan bool),
		wg:       wg,
		hub:      map[string]*websocket.Conn{},
		accepted: map[string]bool{},
		rejected: map[string]bool{},
	}
}

func (c *Web3Connection) Start() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		c.handleWS(w, r)
	})

	c.server.Handler = mux

	c.server.RegisterOnShutdown(func() {
		c.uiEvents.Emit(event_bus.EventLogInfo, "Socket local server stopped")
	})

	go func() {
		for {
			select {
			case event := <-c.w3Events.Events():
				switch event.Type() {
				case event_bus.EventW3WalletAvailable:
					c.handlerWalletAvailable(event.Data())
				case event_bus.EventW3ConnAccepted:
					c.handlerConnAccepted(event.Data())
				case event_bus.EventW3ConnRejected:
					c.handlerConnRejected(event.Data())
				default:
					c.uiEvents.Emit(event_bus.EventLogInfo, fmt.Sprintf(
						"[W3 Connector] Undefined event: %d",
						event.Type()),
					)
				}
			}
		}
	}()

	go func() {
		go func() {
			c.server.ListenAndServe()
			c.uiEvents.Emit(event_bus.EventLogInfo, "[W3 Connector] Stopping local server")
		}()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		<-c.done
		err := c.server.Shutdown(shutdownCtx)
		for connectionId := range c.hub {
			if conn, ok := c.hub[connectionId]; ok {
				_ = conn.Close()
			}
			delete(c.hub, connectionId)
		}

		if err != nil {
			c.uiEvents.Emit(event_bus.EventLogError, fmt.Sprintf("[W3 Connector] Cannot stop server: %s", err))
		}
		return
	}()
	return nil
}

func (c *Web3Connection) handleWS(w http.ResponseWriter, r *http.Request) {
	// Accept only local connections
	if !isRemoteLocal(r.RemoteAddr) {
		w.WriteHeader(403)

		httpResponse := c.newRPCResponse(101, &ResponseErrorFatal{
			Error: "only_local_accepted",
		}).toJSON()

		_, _ = w.Write(httpResponse)
		return
	}

	// Accept only with X-Extension header
	extensionId := r.URL.Query().Get("id")
	if len(extensionId) != 36 {
		w.WriteHeader(400)

		httpResponse := c.newRPCResponse(101, &ResponseErrorFatal{
			Error: "bad_extension_id",
		}).toJSON()

		_, _ = w.Write(httpResponse)
		return
	}

	if c.rejected[extensionId] {
		w.WriteHeader(503)
		httpResponse := c.newRPCResponse(101, &ResponseErrorFatal{
			Error: "rejected",
		}).toJSON()

		_, _ = w.Write(httpResponse)
		return
	}

	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(500)

		httpResponse := c.newRPCResponse(101, &ResponseErrorFatal{
			Error: "upgrader_error",
		}).toJSON()

		_, _ = w.Write(httpResponse)
		return
	}

	defer func() {
		_ = conn.Close()
		delete(c.hub, r.RemoteAddr)
	}()

	c.hub[extensionId] = conn

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			c.uiEvents.Emit(event_bus.EventLogError, fmt.Sprintf("[W3 Connector] Connection error: %s", err))
			break
		}

		//c.uiEvents.Emit(event_bus.EventLogInfo, fmt.Sprintf("[W3 Connector] Message got: %s", message))

		parsedRequest := &RPCMessageReq{}
		err = json.Unmarshal(message, parsedRequest)
		if err != nil {
			c.uiEvents.Emit(event_bus.EventLogWarning, fmt.Sprintf("[W3 Connector] Undefined message: %s", message))
		}

		if parsedRequest.Id != extensionId {
			rpcMessage := c.newRPCResponse(respCodeErrorFatal, &ResponseErrorFatal{
				Error: "not_authorized",
			})
			_ = conn.WriteJSON(rpcMessage)
		}

		switch parsedRequest.Type {
		case reqCodeConnect:
			if c.accepted[extensionId] {
				c.handlerConnAccepted(&dto.ResponseAcceptDTO{
					InstanceId: extensionId,
				})
			} else {
				c.uiEvents.Emit(event_bus.EventW3Connect, &dto.ConnectDTO{
					InstanceId: extensionId,
					Origin:     r.Header.Get("Origin"),
					RemoteAddr: conn.RemoteAddr().String(),
				})
			}
		case reqCodePing:
			c.handlerWalletPing(extensionId)
		case reqCodeRequestAccounts:
			payload := parsedRequest.Data.(*GetAccountsRequest)
			c.uiEvents.Emit(event_bus.EventW3RequestAccounts, &dto.RequestAccountsDTO{
				InstanceId: extensionId,
				Origin:     r.Header.Get("Origin"),
				ChainKey:   payload.ChainKey,
			})
		/*
			case "stop":
				conn.Close()
				c.Stop()
				return*/
		default:
			c.uiEvents.Emit(event_bus.EventLogWarning, fmt.Sprintf("[W3 Connector] Got undefined message: %s", message))
		}
	}
}

func isRemoteLocal(addr string) bool {
	remoteIP, _, err := net.SplitHostPort(addr)

	if err != nil {
		return false
	}

	parsedIP := net.ParseIP(remoteIP)

	if parsedIP == nil || !parsedIP.IsLoopback() {
		return false
	}
	return true
}

func (c *Web3Connection) newRPCResponse(msgType uint16, data interface{}) *RPCMessageResp {
	return &RPCMessageResp{
		RPCMessageHeader: &RPCMessageHeader{
			Version: protocolVersion,
			Id:      c.walletId,
			Type:    msgType,
		},
		Data: data,
	}
}

func (c *Web3Connection) isWalletAvailable() bool {
	return c.walletId != ``
}

func (c *Web3Connection) walletStatus() uint8 {
	if c.isWalletAvailable() {
		return 1 // Available
	} else {
		return 0 // Not available
	}
}

func (c *Web3Connection) Stop() {
	defer c.wg.Done()
	fmt.Println("[Web3] Stopping")
	// c.uiEvents.Emit(event_bus.EventLogInfo, "Trying stopping socket server")
	c.done <- true
}
