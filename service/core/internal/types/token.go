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
	TokenBase    = TokenStandard(`Base`)
	TokenERC20   = TokenStandard(`ERC-20`)
	TokenERC721  = TokenStandard(`ERC-771`)
	TokenERC777  = TokenStandard(`ERC-777`)
	TokenERC1155 = TokenStandard(`ERC-1155`)
	TokenERC4626 = TokenStandard(`ERC-4626`)
	TokenTRC10   = TokenStandard(`TRC-10`)
	TokenTRC20   = TokenStandard(`TRC-20`)
	TokenBEP20   = TokenStandard(`TRC-20`)

	ContractBase = `__base`
	ContractZero = `0x0000000000000000000000000000000000000000`
)

var (
	registeredTokenStandards = map[TokenStandard]bool{
		TokenBase:    true,
		TokenERC20:   true,
		TokenERC721:  true,
		TokenERC777:  true,
		TokenERC1155: true,
		TokenERC4626: true,
		TokenTRC10:   true,
		TokenTRC20:   true,
	}
	activesTokenStandards = map[mhda.NetworkType][]TokenStandard{
		mhda.EthereumVM: {TokenERC20, TokenERC721, TokenERC1155},
		mhda.TronVM:     {TokenTRC20, TokenTRC10},
	}
)

type TokenStandard string

func GetTokenStandardNamesByChain(networkType mhda.NetworkType) []string {
	var result []string
	for _, standardName := range activesTokenStandards[networkType] {
		result = append(result, string(standardName))
	}
	return result
}

func GetTokenStandards(networkType mhda.NetworkType) []TokenStandard {
	return activesTokenStandards[networkType]
}

func GetTokenStandByName(str string) TokenStandard {
	if _, ok := registeredTokenStandards[TokenStandard(str)]; ok {
		return TokenStandard(str)
	} else {
		return ``
	}
}
