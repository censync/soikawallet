package network

import (
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/config/chain"
	"github.com/censync/soikawallet/config/networks"
	"github.com/censync/soikawallet/service/wallet/internal/network/btc"
	"github.com/censync/soikawallet/service/wallet/internal/network/evm"
	"github.com/censync/soikawallet/service/wallet/internal/network/tron"
	"github.com/censync/soikawallet/types"
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
		chain.EthereumChain.Key():     evm.NewEVM(networks.Ethereum),
		chain.TronChain.Key():         tron.NewTron(networks.Tron),
		chain.PolygonChain.Key():      evm.NewEVM(networks.Polygon),
		chain.BinanceSmartChain.Key(): evm.NewEVM(networks.BSC),
		chain.AvalancheCChain.Key():   evm.NewEVM(networks.AvalancheC),
		chain.OptimismChain.Key():     evm.NewEVM(networks.Optimism),
		chain.ArbitrumChain.Key():     evm.NewEVM(networks.ArbitrumOne),
		chain.Moonbeam.Key():          evm.NewEVM(networks.Moonbeam),
		chain.BaseChain.Key():         evm.NewEVM(networks.Base),
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
