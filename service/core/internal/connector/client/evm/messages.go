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
	"encoding/json"
	"errors"
	ws "github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
)

type jsonRPCRawResponse struct {
	Version string `json:"jsonrpc,omitempty"`
	Id      string `json:"id,omitempty"`

	// Call response
	Result interface{} `json:"result,omitempty"`

	// Subscription response
	Method string         `json:"method,omitempty"`
	Params *jsonRPCParams `json:"params,omitempty"`

	Error *jsonRPCError
}

type jsonRPCRawBatchResponse []jsonRPCRawResponse

func (r *jsonRPCRawResponse) isCallResponse() bool {
	return r.Result != nil || r.Error != nil
}

func (r *jsonRPCRawResponse) isSubResponse() bool {
	return r.Method == "eth_subscription" && r.Params != nil
}

type jsonRPCParams struct {
	Subscription string      `json:"subscription,omitempty"`
	Result       interface{} `json:"result,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type responsesPool struct {
	// Mutex requests
	mur sync.RWMutex
	// Mutex batch requests
	mub sync.RWMutex
	// Mutex subscriptions
	mus            sync.RWMutex
	responses      chan *jsonRPCRawResponse
	batchResponses chan *jsonRPCRawBatchResponse

	requests      map[string]*awaitRequest
	batchRequests map[string]*awaitRequest
	subs          map[string]*awaitRequest
	connErr       chan struct{}
	log           *logrus.Entry
}

type awaitRequest struct {
	respCh   chan interface{}
	errCh    chan *jsonRPCError
	cancelCh chan struct{}
}

func newResponsesPool(logger *logrus.Entry) *responsesPool {
	logger = logger.WithFields(logrus.Fields{
		"sub_module": "responses_pool",
	})
	return &responsesPool{
		requests:       map[string]*awaitRequest{},
		batchRequests:  map[string]*awaitRequest{},
		responses:      make(chan *jsonRPCRawResponse),
		batchResponses: make(chan *jsonRPCRawBatchResponse),
		subs:           map[string]*awaitRequest{},
		connErr:        make(chan struct{}),
		log:            logger,
	}
}

func (m *responsesPool) registerRequest(requestId string) {
	m.log.Debugf("Register request %s", requestId)
	m.mur.Lock()
	defer m.mur.Unlock()
	m.requests[requestId] = &awaitRequest{
		respCh:   make(chan interface{}),
		errCh:    make(chan *jsonRPCError),
		cancelCh: make(chan struct{}),
	}
}

func (m *responsesPool) registerBatchRequest(requestId string) {
	m.log.Debugf("Register batch request %s", requestId)
	m.mub.Lock()
	defer m.mub.Unlock()
	m.batchRequests[requestId] = &awaitRequest{
		respCh:   make(chan interface{}),
		errCh:    make(chan *jsonRPCError),
		cancelCh: make(chan struct{}),
	}
}

func (m *responsesPool) unregisterRequest(requestId string) {
	m.log.Debugf("Unregister request %s", requestId)
	m.mur.Lock()
	defer m.mur.Unlock()
	close(m.requests[requestId].respCh)
	close(m.requests[requestId].errCh)
	close(m.requests[requestId].cancelCh)
	delete(m.requests, requestId)
}

func (m *responsesPool) unregisterBatchRequest(requestId string) {
	m.log.Debugf("Unregister batch request %s", requestId)
	m.mub.Lock()
	defer m.mub.Unlock()
	close(m.batchRequests[requestId].respCh)
	close(m.batchRequests[requestId].errCh)
	close(m.batchRequests[requestId].cancelCh)
	delete(m.batchRequests, requestId)
}

func (m *responsesPool) registerSubscription(subId string) (<-chan interface{}, <-chan struct{}) {
	m.log.Debugf("Register subscription %s", subId)
	m.mus.Lock()
	defer m.mus.Unlock()

	respCh := make(chan interface{})
	cancelCh := make(chan struct{})

	m.subs[subId] = &awaitRequest{
		respCh: respCh,
		// errCh:    make(chan *jsonRPCError),
		cancelCh: cancelCh,
	}
	return respCh, cancelCh
}

func (m *responsesPool) unregisterSubscription(subId string) {
	m.log.Debugf("Unregister subscription %s", subId)
	m.mus.Lock()
	defer m.mus.Unlock()
	m.subs[subId].cancelCh <- struct{}{}
	close(m.subs[subId].respCh)
	close(m.subs[subId].cancelCh)
	delete(m.subs, subId)
}

func (c *ClientEVM) connListener() {
	c.log.Debug("Starting conn listener")
	for {
		_, data, err := c.conn.ReadMessage()
		// c.log.Debug(string(data))
		if err != nil {
			// conn error
			if errors.Is(err, ws.ErrCloseSent) {
				c.log.Info("Connection closed gracefully")
			} else {
				var closeMessage *ws.CloseError
				ok := errors.As(err, &closeMessage)
				if ok {
					c.log.Warnf("Connection closed CloseError: %d %s", closeMessage.Code, closeMessage.Text)
				} else {
					c.log.Warnf("Connection closed with unexpected error: %s", err)
				}
			}

			c.isConnected = false
			c.rPool.connErr <- struct{}{} // based on error type, should apply or not to reconnect action
			c.log.Debug("Stopping connListener with conn closing")
			return
		}

		resp := &jsonRPCRawResponse{}

		err = json.Unmarshal(data, &resp)

		if err == nil {
			c.rPool.responses <- resp
		} else {
			respBatch := &jsonRPCRawBatchResponse{}

			err = json.Unmarshal(data, &respBatch)

			if err != nil {
				c.log.Errorf("Connection unmarshal %s", err)
				// cannot unmarshal response
				c.rPool.connErr <- struct{}{}
				c.log.Debug("Stopping connListener with unmarshal error")
				return
			}
			c.rPool.batchResponses <- respBatch
		}
	}
}

func (m *responsesPool) requestHandler() {
	m.log.Debugf("Starting response pool")
	for {
		// Process response
		select {
		case <-m.connErr:
			m.log.Infof("Handled connErr")
			for subId := range m.subs {
				m.unregisterSubscription(subId)
			}
			return
		// m.isReady = false
		// log.Println("Connection closed. Pool size", m.size())
		case response := <-m.responses:
			if response.isCallResponse() {
				// time.Sleep(2 * time.Second)
				if callResponse, ok := m.requests[response.Id]; ok {
					if response.Error == nil {
						callResponse.respCh <- response.Result
					} else {
						callResponse.errCh <- response.Error
					}
				} else {
					m.log.Warnf("Undefined request id: %v", response)
				}
			} else if response.isSubResponse() {
				if subResponse, ok := m.subs[response.Params.Subscription]; ok {
					subResponse.respCh <- response.Params.Result
				} else {
					m.log.Warnf("Subscription undefined request id: %v", response.Params.Subscription)
				}
			} else {
				m.log.Warnf("Undefined response  %v", response)
			}
		case batchResponse := <-m.batchResponses:
			m.log.Infof("Batch response %v", batchResponse)
		}
	}
}
