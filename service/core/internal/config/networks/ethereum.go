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

//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	types2 "github.com/censync/soikawallet/service/core/internal/types"
)

var Ethereum = types2.NewNetwork(
	mhda.ETH,
	`Ethereum`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types2.EVMConfig{
		ChainId: 0x1,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://etherscan.io/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types2.TokenERC20,
	`Tether USD`,
	`USDT`,
	`0xdAC17F958D2ee523a2206206994597C13D831ec7`,
	6,
).SetBuiltinToken(
	types2.TokenERC20,
	`USD network`,
	`USDC`,
	`0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`,
	6,
).SetBuiltinToken(
	types2.TokenERC20,
	`Binance USD`,
	`BUSD`,
	`0x4Fabb145d64652a948d72533023f6E7A623C7C53`,
	18,
).SetBuiltinToken(
	types2.TokenERC20,
	`Matic Token`,
	`MATIC`,
	`0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0`,
	18,
).SetBuiltinToken(
	types2.TokenERC20,
	`Dai Stablecoin`,
	`DAI`,
	`0x6B175474E89094C44Da98b954EedeAC495271d0F`,
	18,
).SetBuiltinToken(
	types2.TokenERC20,
	`ChainLink Token`,
	`LINK`,
	`0x514910771AF9Ca656af840dff83E8264EcF986CA`,
	18,
).SetBuiltinToken(
	types2.TokenERC20,
	`SHIBA INU`,
	`SHIB`,
	`0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE`,
	18,
)
