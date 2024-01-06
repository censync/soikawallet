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

package connector

import (
	"errors"
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/config/version"
	"github.com/censync/soikawallet/service/core/internal/connector/client"
	"github.com/censync/soikawallet/service/core/internal/connector/client/evm"
	"github.com/censync/soikawallet/service/core/internal/connector/config"
	"github.com/censync/soikawallet/service/core/internal/connector/provider"
	evmWrapper "github.com/censync/soikawallet/service/core/internal/connector/provider/evm"
	tronWrapper "github.com/censync/soikawallet/service/core/internal/connector/provider/tron"
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

type Pool map[mhda.ChainKey]map[uint32]client.Client

type Connector struct {
	connPool Pool
	sHub     *sub.Hub
	networks ProviderMap
}

var (
	log                     = logrus.WithFields(logrus.Fields{"service": "connector", "module": "connector"})
	ErrNoAvailableEndpoints = errors.New("no available rpc endpoints")
	defaultHeaders          = http.Header{
		"User-Agent": []string{fmt.Sprintf("Soika Wallet Client %s", version.VERSION)},
	}
)

type ProviderMap map[mhda.ChainKey]provider.Provider

func NewConnector(chains map[mhda.ChainKey][]*types.RPC) (*Connector, error) {

	connector := &Connector{
		connPool: Pool{},
		sHub:     sub.New(),
		networks: ProviderMap{},
	}

	for chainKey, rpcEndpoints := range chains {
		// fix linking
		connector.networks[chainKey] = evmWrapper.NewEVM()

		err := connector.addChain(chainKey, rpcEndpoints...)
		if errors.Is(err, ErrNoAvailableEndpoints) {
			log.Errorf("Cannot add chain :%s", err)
			//return nil, errors.New(fmt.Sprintf("no available rpc endpoints for chain: %s", chainKey))
			continue
		}

		// init base subs to added pool
		if subs, ok := config.Subs[chainKey]; ok {
			err = connector.initSubscriptions(chainKey, subs)
			if err != nil {
				log.Errorf("Cannot init subscriptions for chain: %s", err)
				// return nil, err
			}
		} else {
			log.Infof("No subscriptions for chain: %s", chainKey)
		}

	}

	return connector, nil
}

func (c *Connector) AvailableChains() []mhda.ChainKey {
	return c.GetConnectedChains()
}

func (c *Connector) getClient(chainKey mhda.ChainKey) (client.Client, error) {
	_, ok := c.connPool[chainKey]
	if !ok {
		return nil, fmt.Errorf("pool not initialized for chain key: %s", chainKey)
	}

	availableClients := []uint32{}

	for index := range c.connPool[chainKey] {
		if c.connPool[chainKey][index].IsReady() {
			availableClients = append(availableClients, index)
		}

	}

	readyPoolSize := len(availableClients)

	selectedIndex := availableClients[uint32(rand.Intn(readyPoolSize))]

	log.Debugf("Selected client %d from [%v]", selectedIndex, availableClients)

	return c.connPool[chainKey][selectedIndex], nil
}

func (c *Connector) GetProvider(ctx *types.RPCContext) (provider.Provider, error) {
	adapter, ok := c.networks[ctx.ChainKey()]

	if !ok {
		return nil, errors.New("chain not configured")
	}

	connClient, err := c.getClient(ctx.ChainKey())
	if err != nil {
		return nil, err
	}

	switch adapter.GetType() {
	case evmWrapper.ProviderTypeEVM:
		return adapter.(*evmWrapper.EVM).WithClient(ctx, connClient)
	case tronWrapper.ProviderTypeTron:
		return adapter.(*tronWrapper.Tron).WithClient(ctx, connClient)
	default:
		return nil, errors.New("undefined adapter")
	}
}

func (c *Connector) addChain(chainKey mhda.ChainKey, rpcs ...*types.RPC) error {
	connected := 0

	if _, ok := c.connPool[chainKey]; !ok {
		c.connPool[chainKey] = map[uint32]client.Client{}
	}

	for _, rpc := range rpcs {
		rpcConnection := evm.NewClientEVM(rpc.Index(), rpc.Endpoint(), defaultHeaders)

		for attempt := connectionAttempts; attempt > 0; attempt-- {
			err := rpcConnection.Dial()
			if err == nil {
				log.Debugf("Connected to rpc: %s", rpc.Endpoint())
				c.connPool[chainKey][rpc.Index()] = rpcConnection
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

func (c *Connector) initSubscriptions(chainKey mhda.ChainKey, subs []*sub.Subscription) error {
	if _, ok := c.connPool[chainKey]; ok {
		if len(c.connPool[chainKey]) == 0 {
			return errors.New("no available connections")
		}

		for _, subEntry := range subs {
			// https://tip.golang.org/wiki/LoopvarExperiment
			connClient, err := c.getClient(chainKey)
			if err != nil {
				// add attempts
				continue
			}
			subCh, cancelCh, err := connClient.StartSubscription(subEntry.Method, subEntry.Params)

			if err != nil {
				// log error or cancel and revert
				log.Warnf("Cannot init subscription %s", err)
				continue
			} else {
				// addr of sub: chainKey, nodeIndex, subId
				log.Debugf("Initialized subsctiption %s %d", chainKey, connClient.Index())

				err = c.sHub.Add(chainKey, subEntry, connClient.Index(), subCh, cancelCh)

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

func (c *Connector) GetConnectedChains() []mhda.ChainKey {
	var chainKeys []mhda.ChainKey
	for chainKey := range c.connPool {
		chainKeys = append(chainKeys, chainKey)
	}
	return chainKeys
}

func (c *Connector) Shutdown() {
	for chainKey := range c.connPool {

		for connId := range c.connPool[chainKey] {
			c.connPool[chainKey][connId].Stop()
		}
		// before next steps, subs must be closed, and requests finished
		//p.Unsubscribe(chainKey, "*")

		for connId := range c.connPool[chainKey] {
			c.connPool[chainKey][connId].Disconnect()
			delete(c.connPool[chainKey], connId)
		}
	}
}
