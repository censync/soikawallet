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

package subscriptions

import (
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	Latest      = Type(1) // base network subscriptions (latest block, safe, finalized, slashing, etc)
	Syncing     = Type(2)
	Safe        = Type(3)
	Finalized   = Type(4)
	Slashing    = Type(5)
	EventFilter = Type(6) // base temporary subscriptions (logs with filters)

	Dummy  = Implementation(0) // dummy subscription implementation with loop, if websocket isn't available
	Socket = Implementation(1)

	Idle    = Status(0)
	Blocked = Status(1)
	Ready   = Status(2)
	Closing = Status(3) // if normal or finished

	CloseNormal      = 1
	CloseByFinishing = 2
	CloseByConn      = 3 // CloseAll or Reconnect
	CloseByRpc       = 4
)

type Hub struct {
	subscriptions map[mhda.ChainKey]map[string]*Subscription
	mu            sync.RWMutex
}

type Type uint8

type Implementation uint8

type Status uint8

type Subscription struct {
	Name    string
	Impl    Implementation
	SubType Type
	Status  Status
	Method  string
	Params  []interface{}
	Func    func(interface{})
}

type Event struct {
	ChainKey  mhda.ChainKey
	NodeIndex uint32
	Data      interface{}
}

var (
	log = logrus.WithFields(logrus.Fields{"service": "connector", "module": "subscription"})
)

func New() *Hub {
	return &Hub{
		subscriptions: map[mhda.ChainKey]map[string]*Subscription{},
	}
}

func (h *Hub) Add(chainKey mhda.ChainKey, sub *Subscription, index uint32, respCh <-chan interface{}, cancelCh <-chan struct{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.subscriptions[chainKey]; !ok {
		h.subscriptions[chainKey] = map[string]*Subscription{}
	}
	if _, ok := h.subscriptions[chainKey][sub.Name]; ok {
		return errors.New("subscription with this name already registered")
	}
	sub.Status = Ready
	h.subscriptions[chainKey][sub.Name] = sub

	go func() {
		for {
			select {
			case <-cancelCh:
				log.Infof("Unregister subscription from hub by conn: %s %d", chainKey, index)
				// unregister ?
				return
			case event, ok := <-respCh:
				if !ok {
					log.Errorf("HUB: Handler got event on closed channel %s %d", chainKey, index)
				}
				sub.Func(event)
			}
		}
	}()
	return nil
}

func (h *Hub) unregister(chainKey mhda.ChainKey, name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscriptions[chainKey], name)
}
