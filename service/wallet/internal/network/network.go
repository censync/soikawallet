package network

import (
	"errors"
	"github.com/censync/soikawallet/config/networks"
	"github.com/censync/soikawallet/service/wallet/internal/network/evm"
	"github.com/censync/soikawallet/service/wallet/internal/network/tron"
	"github.com/censync/soikawallet/types"
	"sync"
)

type Provider struct {
	mu       sync.RWMutex
	networks map[types.CoinType]types.NetworkAdapter
}

var networkProviders = &Provider{
	networks: map[types.CoinType]types.NetworkAdapter{
		types.Ethereum: evm.NewEVM(networks.Ethereum),
		types.Tron:     tron.NewTron(networks.Tron),
		types.Polygon:  evm.NewEVM(networks.Polygon),
	},
}

func (s *Provider) IsExists(index types.CoinType) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.networks[index]
	return ok
}

func Get(index types.CoinType) types.NetworkAdapter {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	if network, ok := networkProviders.networks[index]; ok {
		return network
	}
	return nil
}

func WithContext(ctx *types.RPCContext) (types.NetworkAdapter, error) {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	if !types.IsCoinExists(ctx.CoinType()) {
		return nil, errors.New("coin type is not set")
	}

	network, ok := networkProviders.networks[ctx.CoinType()]

	if !ok {
		return nil, errors.New("network is not defined")
	}
	return network, nil
}
