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
	"encoding/json"
	"github.com/sirupsen/logrus"
	"sync"
)

type jsonRPCRawResponse struct {
	Version string `json:"jsonrpc,omitempty"`
	Id      string `json:"id,omitempty"`

	// Call response
	Result json.RawMessage `json:"result,omitempty"`

	// Subscription response
	Method string         `json:"method,omitempty"`
	Params *jsonRPCParams `json:"params,omitempty"`

	Error *jsonRPCError
}

func (r *jsonRPCRawResponse) isCallResponse() bool {
	return r.Result != nil || r.Error != nil
}

func (r *jsonRPCRawResponse) isSubResponse() bool {
	return r.Method == "eth_subscription" && r.Params != nil
}

type jsonRPCRawBatchResponse []jsonRPCRawResponse

func (br *jsonRPCRawBatchResponse) IsEmpty() bool {
	return br == nil || len(*br) > 0
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
	// Mutex subscriptions
	mus            sync.RWMutex
	responses      chan *jsonRPCRawResponse
	batchResponses chan *jsonRPCRawBatchResponse

	requests      *awaitRequestsPool
	batchRequests *awaitRequestsPool
	subs          map[string]*awaitRequest
	connErr       chan struct{}
	log           *logrus.Entry
}

func newResponsesPool(logger *logrus.Entry) *responsesPool {
	logger = logger.WithFields(logrus.Fields{
		"sub_module": "responses_pool",
	})
	return &responsesPool{
		requests:       newAwaitRequestsPool(),
		batchRequests:  newAwaitRequestsPool(),
		responses:      make(chan *jsonRPCRawResponse),
		batchResponses: make(chan *jsonRPCRawBatchResponse),
		subs:           map[string]*awaitRequest{},
		connErr:        make(chan struct{}),
		log:            logger,
	}
}

type awaitRequest struct {
	// respCh - responses handling channel for ws request
	respCh chan interface{}
	// errCh - errors handling channel for ws request
	errCh chan *jsonRPCError
	// cancelCh - cancellation operations handling channel for request
	cancelCh chan struct{}
}

type awaitRequestsPool struct {
	mu       sync.Mutex
	requests map[string]*awaitRequest
}

func (rp *awaitRequestsPool) Get(requestId string) *awaitRequest {
	return rp.requests[requestId]
}

func (rp *awaitRequestsPool) SetRequests(requests map[string]*awaitRequest) {
	rp.requests = requests
}

func newAwaitRequestsPool() *awaitRequestsPool {
	return &awaitRequestsPool{requests: map[string]*awaitRequest{}}
}

func (rp *awaitRequestsPool) register(requestId string) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	rp.requests[requestId] = &awaitRequest{
		respCh:   make(chan interface{}),
		errCh:    make(chan *jsonRPCError),
		cancelCh: make(chan struct{}),
	}
}

func (rp *awaitRequestsPool) unregister(requestId string) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	close(rp.requests[requestId].respCh)
	close(rp.requests[requestId].errCh)
	close(rp.requests[requestId].cancelCh)
	delete(rp.requests, requestId)
}

// Call await requests helpers
func (m *responsesPool) registerRequest(requestId string) {
	m.log.Debugf("Register request %s", requestId)
	m.requests.register(requestId)
}

func (m *responsesPool) unregisterRequest(requestId string) {
	m.log.Debugf("Unregister request %s", requestId)
	m.requests.unregister(requestId)
}

// Batch call await requests helpers
func (m *responsesPool) registerBatchRequest(requestId string) {
	m.log.Debugf("Register batch request %s", requestId)
	m.batchRequests.register(requestId)
}

func (m *responsesPool) unregisterBatchRequest(requestId string) {
	m.log.Debugf("Unregister batch request %s", requestId)
	m.batchRequests.unregister(requestId)
}

// Subscription await requests helpers
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
				if registeredResponse := m.requests.Get(response.Id); registeredResponse != nil {
					if response.Error == nil {
						registeredResponse.respCh <- response.Result
					} else {
						registeredResponse.errCh <- response.Error
					}
				} else {
					m.log.Warnf("Undefined call request id: %v", response)
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

			if !batchResponse.IsEmpty() {
				reqIdSum := bytes.NewBuffer([]byte{})

				for index := range *batchResponse {
					reqIdSum.WriteString((*batchResponse)[index].Id)
				}

				reqIdSumHash := md5.Sum(reqIdSum.Bytes())
				reqIdSumHashStr := hex.EncodeToString(reqIdSumHash[:])

				if registeredBatchResponse := m.batchRequests.Get(reqIdSumHashStr); registeredBatchResponse != nil {
					registeredBatchResponse.respCh <- batchResponse
				} else {
					m.log.Warnf("Undefined batch request id: %v", reqIdSumHashStr)
				}

			} else {
				m.log.Warnf("Batch response is empty")
			}
		}
	}
}
