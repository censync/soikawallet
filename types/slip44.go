package types

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/config/chain"

	"sort"
)

var (
	registeredNetworks = map[string]*mhda.Chain{
		`Bitcoin`:     chain.BitcoinChain,
		`Ethereum`:    chain.EthereumChain,
		`Tron`:        chain.TronChain,
		`Polygon`:     chain.PolygonChain,
		`BSC`:         chain.BinanceSmartChain,
		`Optimism`:    chain.OptimismChain,
		`Arbitrum`:    chain.ArbitrumChain,
		`Avalanche C`: chain.AvalancheCChain,
		`Moonbeam`:    chain.Moonbeam,
		`Base`:        chain.BaseChain,
	}

	registeredNetworksIndexes    = map[*mhda.Chain]string{}
	registeredNetworksTypes      []mhda.ChainKey
	registeredNetworksNames      []string
	registeredNetworksNamesIndex = map[mhda.ChainKey]string{}
)

func init() {
	for name, chainKey := range registeredNetworks {
		registeredNetworksIndexes[chainKey] = name
		registeredNetworksTypes = append(registeredNetworksTypes, chainKey.Key())
		registeredNetworksNames = append(registeredNetworksNames, name)
		registeredNetworksNamesIndex[chainKey.Key()] = name
	}
	sort.Strings(registeredNetworksNames)
}

func GetChainNames() []string {
	return registeredNetworksNames
}

func GetChains() []mhda.ChainKey {
	return registeredNetworksTypes
}

func GetChainByName(str string) *mhda.Chain {
	if networkType, ok := registeredNetworks[str]; ok {
		return networkType
	} else {
		return nil
	}
}

func GetNetworkNameByKey(chainKey mhda.ChainKey) string {
	return registeredNetworksNamesIndex[chainKey]
}

func IsNetworkExists(val mhda.ChainKey) bool {
	// TODO: Fix
	/*for i := range registeredNetworksTypes {
		if registeredNetworksTypes[i] == val {
			return true
		}
	}*/
	//return false
	return true
}
