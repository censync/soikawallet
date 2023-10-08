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

import mhda "github.com/censync/go-mhda"

const (
	TokenBase    = TokenStandard(1)
	TokenERC20   = TokenStandard(20)
	TokenERC721  = TokenStandard(21)
	TokenERC777  = TokenStandard(22)
	TokenERC1155 = TokenStandard(23)
	TokenERC4626 = TokenStandard(24)
	TokenTRC10   = TokenStandard(40)
	TokenTRC20   = TokenStandard(41)
	TokenBEP20   = TokenStandard(50)

	ContractBase = `__base`
	ContractZero = `0x0000000000000000000000000000000000000000`
)

var (
	registeredTokenStandards = map[TokenStandard]string{
		TokenBase:    `Base`,
		TokenERC20:   `ERC-20`,
		TokenERC721:  `ERC-771`,
		TokenERC777:  `ERC-777`,
		TokenERC1155: `ERC-1155`,
		TokenERC4626: `ERC-4626`,
		TokenTRC10:   `TRC-10`,
		TokenTRC20:   `TRC-20`,
	}
	activesTokenStandards = map[mhda.NetworkType][]TokenStandard{
		mhda.EthereumVM: {TokenERC20, TokenERC721, TokenERC1155},
		mhda.TronVM:     {TokenTRC20, TokenTRC10},
	}
	registeredTokenStandardNames = map[mhda.NetworkType][]string{}
	registeredTokenIndexes       = map[string]TokenStandard{}
)

type TokenStandard uint8

func init() {
	for tokenStandard, tokenStandardName := range registeredTokenStandards {
		registeredTokenIndexes[tokenStandardName] = tokenStandard
	}

	for networkType, activeTokenStandards := range activesTokenStandards {
		names := make([]string, 0)
		for _, tokenStandard := range activeTokenStandards {
			names = append(names, registeredTokenStandards[tokenStandard])
		}
		//sort.Strings(names)
		registeredTokenStandardNames[networkType] = names
	}
}

func GetTokenStandardNames(networkType mhda.NetworkType) []string {
	return registeredTokenStandardNames[networkType]
}

func GetTokenStandardNamesByChain(networkType mhda.ChainKey) []string {
	// mhda.ParseURN()
	//return registeredTokenStandardNames[networkType]
	// TODO: Remove debug
	return []string{`ERC-20`, `ERC-771`, `ERC-1155`}
}

func GetTokenStandards(networkType mhda.NetworkType) []TokenStandard {
	return activesTokenStandards[networkType]
}

func GetTokenStandByName(str string) TokenStandard {
	if tokenStandard, ok := registeredTokenIndexes[str]; ok {
		return tokenStandard
	} else {
		return 0
	}
}
