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

package pool

import (
	"errors"
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/config/version"
	"github.com/censync/soikawallet/service/core/internal/connector/client"
	"github.com/censync/soikawallet/service/core/internal/connector/client/evm"
	sub "github.com/censync/soikawallet/service/core/internal/connector/subscriptions"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const (
	connectionAttempts = 5
	reconnectTimeout   = 500
)

var (
	log                     = logrus.WithFields(logrus.Fields{"service": "connector", "module": "pool"})
	ErrNoAvailableEndpoints = errors.New("no available rpc endpoints")
	defaultHeaders          = http.Header{
		"User-Agent": []string{fmt.Sprintf("Soika Wallet Client %s", version.VERSION)},
	}
)

type Pool struct {
	pool map[mhda.ChainKey]map[uint32]client.Client
	sHub *sub.Hub
}

func NewPool() *Pool {
	hub := sub.New()
	return &Pool{
		pool: map[mhda.ChainKey]map[uint32]client.Client{},
		sHub: hub,
	}
}

func (p *Pool) AddChain(chainKey mhda.ChainKey, rpcs ...*types.RPC) error {
	connected := 0

	if _, ok := p.pool[chainKey]; !ok {
		p.pool[chainKey] = map[uint32]client.Client{}
	}

	for _, rpc := range rpcs {
		rpcConnection := evm.NewClientEVM(rpc.Index(), rpc.Endpoint(), defaultHeaders)

		for attempt := connectionAttempts; attempt > 0; attempt-- {
			err := rpcConnection.Dial()
			if err == nil {
				log.Debugf("Connected to rpc: %s", rpc.Endpoint())
				p.pool[chainKey][rpc.Index()] = rpcConnection
				connected++
				break
			}
			log.Warnf("Cannot connect to rpc (attempt=%d): %s %s", attempt, rpc.Endpoint(), err)
			time.Sleep(reconnectTimeout * time.Millisecond)
		}
	}
	if connected == 0 {
		return ErrNoAvailableEndpoints
	}
	return nil
}

func (p *Pool) InitSubscriptions(chainKey mhda.ChainKey, subs []*sub.Subscription) error {
	if _, ok := p.pool[chainKey]; ok {
		if len(p.pool[chainKey]) == 0 {
			return errors.New("no available connections")
		}

		for _, subEntry := range subs {
			// https://tip.golang.org/wiki/LoopvarExperiment
			client, err := p.GetClient(chainKey)
			if err != nil {
				// add attempts
				continue
			}
			subCh, cancelCh, err := client.StartSubscription(subEntry.Method, subEntry.Params)

			if err != nil {
				// log error or cancel and revert
				log.Warnf("Cannot init subscription %s", err)
				continue
			} else {
				// addr of sub: chainKey, nodeIndex, subId
				log.Debugf("Initialized subsctiption %s %d", chainKey, client.Index())

				err = p.sHub.Add(chainKey, subEntry, client.Index(), subCh, cancelCh)

				if err != nil {
					// cancel subscription connect
					// attempt?
					log.Errorf("Cannot add subscription %s: %s %s", chainKey, subEntry.Name, err)
				} else {
					log.Infof("Subscription added to hub %s: %s %s", chainKey, subEntry.Name, err)
				}
			}
		}

	} else {
		return fmt.Errorf("pool not initialized for chain key: %s", chainKey)
	}
	return nil
}

/*func (p *Pool) Unsubscribe(chainKey mhda.ChainKey, name string) {
	// Change management operations  to Subscription struct
	if _, ok := p.pool[chainKey]; ok {
		if len(p.pool[chainKey]) == 0 {
			log.Debugf("no available connections")
			return
		}

		for nodeIndex := range p.pool[chainKey] {
			log.Debugf("Trying canceling subscriptions from Hub: %s %d", chainKey, nodeIndex)
			// p.pool[chainKey][nodeIndex].CancelSubscriptions()
		}
	}
}*/

func (p *Pool) ConnectedChains() []mhda.ChainKey {
	var chainKeys []mhda.ChainKey
	for chainKey := range p.pool {
		chainKeys = append(chainKeys, chainKey)
	}
	return chainKeys
}

func (p *Pool) GetClient(chainKey mhda.ChainKey) (client.Client, error) {
	_, ok := p.pool[chainKey]
	if !ok {
		return nil, fmt.Errorf("pool not initialized for chain key: %s", chainKey)
	}

	availableClients := []uint32{}

	for index := range p.pool[chainKey] {
		if p.pool[chainKey][index].IsReady() {
			availableClients = append(availableClients, index)
		}

	}

	readyPoolSize := len(availableClients)

	selectedIndex := availableClients[uint32(rand.Intn(readyPoolSize))]

	log.Debugf("Selected client %d from [%v]", selectedIndex, availableClients)

	return p.pool[chainKey][selectedIndex], nil

}

func (p *Pool) StopClient(chainKey mhda.ChainKey) {
	for index := range p.pool[chainKey] {
		p.pool[chainKey][index].Stop()
	}
}

func (p *Pool) CloseAll() {
	for chainKey := range p.pool {

		p.StopClient(chainKey)
		// before next steps, subs must be closed, and requests finished
		//p.Unsubscribe(chainKey, "*")

		for connId := range p.pool[chainKey] {
			p.pool[chainKey][connId].Disconnect()
			delete(p.pool[chainKey], connId)
		}
	}
}
