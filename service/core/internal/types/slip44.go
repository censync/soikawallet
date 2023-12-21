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

package types

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/config/chain"
	"sort"
)

var (
	registeredNetworks = map[string]*mhda.Chain{
		`Ethereum`:    chain.EthereumChain,
		`Polygon`:     chain.PolygonChain,
		`BSC`:         chain.BinanceSmartChain,
		`Avalanche C`: chain.AvalancheCChain,
		`Moonbeam`:    chain.Moonbeam,
		`Bitcoin`:     chain.BitcoinChain,
		`Tron`:        chain.TronChain,

		//  TODO: Add gas calculator for L2 chains
		// `Mantle`: 	  chain.MantleChain,
		// `Optimism`:    chain.OptimismChain,
		// `Arbitrum`:    chain.ArbitrumChain,
		// `Base`:        chain.BaseChain,

		// ZK
		// `zkEVM`:  chain.ZkPolygon,
		// `zkSync`: chain.ZkSyncEra,
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

func GetChainNameByKey(chainKey mhda.ChainKey) string {
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
