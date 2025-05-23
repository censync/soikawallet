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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package clients

import (
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/config/chain"
	"github.com/censync/soikawallet/service/core/internal/clients/btc"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_arbitrum"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_base"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_legacy"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_optimism"
	"github.com/censync/soikawallet/service/core/internal/clients/tron"
	"github.com/censync/soikawallet/service/core/internal/config/networks"
	"github.com/censync/soikawallet/service/core/internal/types"
	"sync"
)

type Provider struct {
	mu              sync.RWMutex
	networks        map[mhda.ChainKey]types.NetworkAdapter
	defaultCurrency string
}

var (
	errNetworkTypeNotSet = errors.New("network type is not set")
	errNetworkNotDefined = errors.New("network is not defined")
)

/*
var (
	_ types.NetworkAdapter = tron.Tron{}
	_ types.NetworkAdapter = evm.EVM{}
	_ types.NetworkAdapter = btc.Bitcoin{}
)*/

var networkProviders = &Provider{
	networks: map[mhda.ChainKey]types.NetworkAdapter{
		chain.BitcoinChain.Key():      btc.NewBTC(networks.Bitcoin),
		chain.EthereumChain.Key():     evm_base.NewEVM(networks.Ethereum),
		chain.TronChain.Key():         tron.NewTron(networks.Tron),
		chain.PolygonChain.Key():      evm_base.NewEVM(networks.Polygon),
		chain.BinanceSmartChain.Key(): evm_legacy.NewEVMLegacy(networks.BSC),
		chain.AvalancheCChain.Key():   evm_base.NewEVM(networks.AvalancheC),
		chain.OptimismChain.Key():     evm_optimism.NewEVMOptimism(networks.Optimism),
		chain.ArbitrumChain.Key():     evm_arbitrum.NewEVMArbitrum(networks.ArbitrumOne),
		chain.Moonbeam.Key():          evm_base.NewEVM(networks.Moonbeam),
		chain.BaseChain.Key():         evm_base.NewEVM(networks.Base),
		chain.ZkPolygon.Key():         evm_base.NewEVM(networks.ZkEVMPolygon),
		chain.ZkSyncEra.Key():         evm_base.NewEVM(networks.ZkSyncEra),
		chain.Linea.Key():             evm_base.NewEVM(networks.Linea),
		chain.MantleChain.Key():       evm_base.NewEVM(networks.Mantle),
		chain.Blast.Key():             evm_base.NewEVM(networks.Blast),
		chain.Gnosis.Key():            evm_base.NewEVM(networks.Gnosis),
	},
	defaultCurrency: `USD`,
}

func (s *Provider) IsExists(chainKey mhda.ChainKey) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.networks[chainKey]
	return ok
}

func Get(chainKey mhda.ChainKey) types.NetworkAdapter {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	if network, ok := networkProviders.networks[chainKey]; ok {
		return network
	}
	return nil
}

func WithContext(ctx *types.RPCContext) (types.NetworkAdapter, error) {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	if !types.IsNetworkExists(ctx.ChainKey()) {
		return nil, errNetworkTypeNotSet
	}

	network, ok := networkProviders.networks[ctx.ChainKey()]

	if !ok {
		return nil, errNetworkNotDefined
	}
	return network, nil
}

func GetAll() map[mhda.ChainKey]types.NetworkAdapter {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	return networkProviders.networks
}
