// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/service/core/internal/connector/client/metrics"
	"github.com/censync/soikawallet/service/core/internal/connector/types/callopts"
	"github.com/ethereum/go-ethereum/log"
	ws "github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"math/rand"
	"net/http"
	"time"
)

const maxBatchParamsCount = 20

type ClientEVM struct {
	isWs        bool
	index       uint32
	rpc         string
	headers     http.Header
	conn        *ws.Conn
	reqs        map[string]chan string
	r           *rand.Rand
	rPool       *responsesPool
	metrics     *metrics.Metrics
	isConnected bool
	stopping    bool
	log         *logrus.Entry
}

type jsonRPCRequest struct {
	Version string        `json:"jsonrpc,omitempty"`
	Id      string        `json:"id,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  []interface{} `json:"params,omitempty"`
}

func NewClientEVM(index uint32, rpc string, headers http.Header) *ClientEVM {
	logger := logrus.WithFields(logrus.Fields{
		"service": "connector",
		"module":  "evm client",
		"rpc":     rpc,
	})

	return &ClientEVM{
		index:   index,
		rpc:     rpc,
		headers: headers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
		metrics: &metrics.Metrics{},
		rPool:   newResponsesPool(logger),
		log:     logger,
	}
}

func (c *ClientEVM) Dial() error {
	log.Debug("Trying dial")

	conn, _, err := ws.DefaultDialer.Dial(c.rpc, c.headers)

	if err != nil {
		c.log.Warnf("Cannot dial to %s", err)
		return err
	}

	c.conn = conn
	go c.rPool.requestHandler()
	go c.connListener()

	// TODO: Check eth_chainId
	c.isConnected = true
	return nil
}

func (c *ClientEVM) Index() uint32 {
	return c.index
}

func (c *ClientEVM) CallBatch(ctx context.Context, opts []*callopts.CallOpts) (interface{}, error) {
	paramsCount := len(opts)

	if paramsCount == 0 {
		return nil, errors.New("empty batch params")
	}

	if paramsCount > maxBatchParamsCount {
		return nil, errors.New("to many batch params")
	}

	batchReqRPC := make([]*jsonRPCRequest, paramsCount)

	reqIdSum := bytes.NewBuffer([]byte{})

	for i := 0; i < paramsCount; i++ {
		reqId := fmt.Sprintf("0x%x", c.r.Uint64())
		reqIdSum.WriteString(reqId)

		batchReqRPC = append(batchReqRPC, &jsonRPCRequest{
			Version: "2.0",
			Id:      reqId,
			Method:  opts[i].Method(),
			Params:  opts[i].Params(),
		})
	}

	reqIdSumHash := md5.Sum(reqIdSum.Bytes())
	reqIdSumHashStr := hex.EncodeToString(reqIdSumHash[:])

	c.rPool.registerBatchRequest(reqIdSumHashStr)
	defer c.rPool.unregisterBatchRequest(reqIdSumHashStr)

	reqIdSum.Reset()

	err := c.conn.WriteJSON(batchReqRPC)
	if err != nil {
		// <-  failed rpc resp
		c.rPool.unregisterBatchRequest(reqIdSumHashStr)
		c.rPool.connErr <- struct{}{}
		return nil, err
	}

	select {
	// Cancelled
	case <-c.rPool.batchRequests[reqIdSumHashStr].cancelCh:
		c.log.Warnf("Conn error %s", reqIdSumHashStr)
		return nil, errors.New("conn error")
	// Deadline
	case <-ctx.Done():
		c.log.Infof("CallBatch finished with deadline %s", reqIdSumHashStr)
		return nil, errors.New("deadline")
	// Response
	case result := <-c.rPool.batchRequests[reqIdSumHashStr].respCh:
		c.log.Infof("CallBatch response received: %s, %s", reqIdSumHashStr, result)
		return result, nil
	// Response (error)
	case resultErr := <-c.rPool.batchRequests[reqIdSumHashStr].errCh:
		c.log.Infof("CallBatch response with error received: %s, %s", reqIdSumHashStr, resultErr)
		return resultErr, nil
	}
}

func (c *ClientEVM) Call(ctx context.Context, method string, params []interface{}) (interface{}, error) {
	return c.CallOpts(ctx, callopts.NewCallOpts(method, params))
}

func (c *ClientEVM) CallOpts(ctx context.Context, opts *callopts.CallOpts) (interface{}, error) {
	if c.stopping && opts.Method() != "eth_unsubscribe" {
		return nil, errors.New("client stopping")
	}

	c.log.Debugf("Starting call method %s [%v]", opts.Method(), opts.Params())

	reqId := fmt.Sprintf("0x%x", c.r.Uint64())
	reqRPC := &jsonRPCRequest{
		Version: "2.0",
		Id:      reqId,
		Method:  opts.Method(),
		Params:  opts.Params(),
	}

	c.rPool.registerRequest(reqId)
	defer c.rPool.unregisterRequest(reqId)

	err := c.conn.WriteJSON(reqRPC)
	if err != nil {
		// <-  failed rpc resp
		c.rPool.unregisterRequest(reqId)
		c.rPool.connErr <- struct{}{}
		return nil, err
	}

	select {
	// Cancelled
	case <-c.rPool.requests[reqId].cancelCh:
		c.log.Warnf("Conn error %s", reqId)
		return nil, errors.New("conn error")
	// Deadline
	case <-ctx.Done():
		c.log.Infof("Call finished with deadline %s", reqId)
		return nil, errors.New("deadline")
	// Response
	case result := <-c.rPool.requests[reqId].respCh:
		c.log.Infof("Response received: %s, %s", reqId, result)
		return result, nil
	// Response (error)
	case resultErr := <-c.rPool.requests[reqId].errCh:
		c.log.Infof("Response with error received: %s, %s", reqId, resultErr)
		return resultErr, nil
	}
}

func (c *ClientEVM) StartSubscription(method string, params []interface{}) (<-chan interface{}, <-chan struct{}, error) {
	if c.stopping {
		return nil, nil, errors.New("client stopping")
	}
	c.log.Debugf("StartSubscription %s %v", method, params)

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(2000*time.Millisecond))

	resp, err := c.Call(ctx, method, params)

	if err != nil {
		return nil, nil, err
	}

	subId, ok := resp.(string)

	if !ok {
		return nil, nil, errors.New("cannot parse subscription id ")
	}

	if subId == "" {
		return nil, nil, errors.New("subscription id is empty")
	}

	respCh, cancelCh := c.rPool.registerSubscription(subId)
	return respCh, cancelCh, nil
}

func (c *ClientEVM) cancelSubscriptions() {
	for subId := range c.rPool.subs {
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1000*time.Millisecond))

		resp, err := c.Call(ctx, "eth_unsubscribe", []interface{}{subId})

		if err != nil {
			c.log.Warnf("Cannot unsubscribe sub: %s", err)
		}

		c.log.Debugf("Subscription succesfully unsubscribed: %s", resp)

		c.rPool.unregisterSubscription(subId)
	}
}

func (c *ClientEVM) Stop() {
	c.stopping = true
	c.cancelSubscriptions()
}

func (c *ClientEVM) IsReady() bool {
	return c.isConnected
}

func (c *ClientEVM) Disconnect() {
	// pools must be awaited and finalized before closure sent
	if !c.stopping {
		c.Stop()
	}
	c.log.Debug("Send normal closure")
	_ = c.conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, ""))
	//c.conn.Close()
}
