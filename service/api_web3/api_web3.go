package api_web3

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

const addr = "127.0.0.1:8114"

type Web3Connection struct {
	events *event_bus.EventBus
	server *http.Server

	upgrader websocket.Upgrader // use default options
	done     chan bool
	wg       *sync.WaitGroup
	hub      map[string]*websocket.Conn
}

func NewWeb3Connection(cfg *config.Config, wg *sync.WaitGroup, events *event_bus.EventBus) *Web3Connection {
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
		events:   events,
		server:   server,
		upgrader: upgrader,
		done:     make(chan bool),
		wg:       wg,
		hub:      map[string]*websocket.Conn{},
	}
}

func (c *Web3Connection) Start() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c.handleWS(w, r)
	})

	c.server.Handler = mux

	c.server.RegisterOnShutdown(func() {
		c.events.Emit(event_bus.EventLogInfo, "Socket local server stopped")
	})

	go func() {
		//defer c.wg.Done()

		go func() {
			c.server.ListenAndServe()
			c.events.Emit(event_bus.EventLogInfo, "Socket stopping local server")
		}()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		<-c.done
		err := c.server.Shutdown(shutdownCtx)
		for id := range c.hub {
			c.hub[id].Close()
		}

		if err != nil {
			c.events.Emit(event_bus.EventLogError, fmt.Sprintf("Socket cannot stop server: %s", err))
		}

	}()
	return nil
}

func (c *Web3Connection) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		conn.Close()
		return
	}

	// Accept only local connections
	if !isRemoteLocal(r.RemoteAddr) {
		conn.Close()
		return
	}

	defer func() {
		conn.Close()
		delete(c.hub, r.RemoteAddr)
	}()

	c.hub[r.RemoteAddr] = conn

	for {
		mt, message, err := conn.ReadMessage()

		if err != nil {
			c.events.Emit(event_bus.EventLogError, fmt.Sprintf("Socket conn err: %s", err))
			break
		}

		c.events.Emit(event_bus.EventLogInfo, fmt.Sprintf("Socket message got: %s", message))

		switch string(message) {
		case "ping":
			c.events.Emit(event_bus.EventLogInfo, "Echo from client")
			err = conn.WriteMessage(mt, []byte("pong"))
			if err != nil {
				c.events.Emit(event_bus.EventLogError, fmt.Sprintf("Socket write err: %s", err))
				break
			}
		case "stop":
			conn.Close()
			c.Stop()
			return
		default:
			c.events.Emit(event_bus.EventLogWarning, fmt.Sprintf("Undefined message: %s", message))
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

func (c *Web3Connection) Stop() {
	defer c.wg.Done()
	c.events.Emit(event_bus.EventLogInfo, "Trying stopping socket server")
	c.done <- true
}
