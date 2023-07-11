package network

import (
	"errors"
	"github.com/censync/soikawallet/config/networks"
	"github.com/censync/soikawallet/service/wallet/internal/network/btc"
	"github.com/censync/soikawallet/service/wallet/internal/network/evm"
	"github.com/censync/soikawallet/service/wallet/internal/network/tron"
	"github.com/censync/soikawallet/types"
	"sync"
)

type Provider struct {
	mu              sync.RWMutex
	networks        map[types.NetworkType]types.NetworkAdapter
	defaultCurrency string
}

var networkProviders = &Provider{
	networks: map[types.NetworkType]types.NetworkAdapter{
		types.Bitcoin:  btc.NewBTC(networks.Bitcoin),
		types.Ethereum: evm.NewEVM(networks.Ethereum),
		types.Tron:     tron.NewTron(networks.Tron),
		types.Polygon:  evm.NewEVM(networks.Polygon),
		types.BSC:      evm.NewEVM(networks.BSC),
	},
	defaultCurrency: `USD`,
}

func (s *Provider) IsExists(index types.NetworkType) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.networks[index]
	return ok
}

func Get(index types.NetworkType) types.NetworkAdapter {
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

	if !types.IsNetworkExists(ctx.Network()) {
		return nil, errors.New("network type is not set")
	}

	network, ok := networkProviders.networks[ctx.Network()]

	if !ok {
		return nil, errors.New("network is not defined")
	}
	return network, nil
}

func GetAll() map[types.NetworkType]types.NetworkAdapter {
	networkProviders.mu.RLock()
	defer networkProviders.mu.RUnlock()

	return networkProviders.networks
}
