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
	"context"
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/connector/config"
	"github.com/censync/soikawallet/service/core/internal/connector/pool"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/sirupsen/logrus"
)

type Connector struct {
	*pool.Pool
}

var (
	log = logrus.WithFields(logrus.Fields{"service": "connector", "module": "connector"})
)

func NewConnector(chains map[mhda.ChainKey][]*types.RPC) (*Connector, error) {
	connPool := pool.NewPool()

	for chainKey, rpcEndpoints := range chains {
		err := connPool.AddChain(chainKey, rpcEndpoints...)
		if errors.Is(err, pool.ErrNoAvailableEndpoints) {
			log.Errorf("Cannot add chain :%s", err)
			//return nil, errors.New(fmt.Sprintf("no available rpc endpoints for chain: %s", chainKey))
			continue
		}

		// init base subs to added pool

		if subs, ok := config.Subs[chainKey]; ok {
			err = connPool.InitSubscriptions(chainKey, subs)
			if err != nil {
				log.Errorf("Cannot init subscriptions for chain: %s", err)
				// return nil, err
			}
		} else {
			log.Infof("No subscriptions for chain: %s", chainKey)
		}

	}

	return &Connector{Pool: connPool}, nil
}

func (c *Connector) PoolChains() []mhda.ChainKey {
	return c.ConnectedChains()
}

func (c *Connector) Call(ctx context.Context, chainKey mhda.ChainKey, method string, params []interface{}) (interface{}, error) {
	client, err := c.GetClient(chainKey)
	if err != nil {
		return nil, err
	}

	return client.Call(ctx, method, params)
}

func (c *Connector) Shutdown() {
	c.CloseAll()
}
