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

//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	types2 "github.com/censync/soikawallet/service/core/internal/types"
)

var Polygon = types2.NewNetwork(
	mhda.MATIC,
	`Polygon`,
	`MATIC`,
	18,
	1e9,
	"gwei",
	true,
	&types2.EVMConfig{
		ChainId: 0x89,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon`, // https://polygon-rpc.com
	`https://polygonscan.com/`,                 // /block/ /address/ /tx/
).SetBuiltinToken(
	types2.TokenERC20,
	`(PoS) Tether USD`,
	`USDT`,
	`0xc2132d05d31c914a87c6611c10748aeb04b58e8f`,
	6,
).SetBuiltinToken(
	types2.TokenERC20,
	`USD network (PoS)`,
	`USDC`,
	`0x2791bca1f2de4661ed88a30c99a7a9449aa84174`,
	6,
).SetBuiltinToken(
	types2.TokenERC20,
	`(PoS) Binance USD`,
	`BUSD`,
	`0xdAb529f40E671A1D4bF91361c21bf9f0C9712ab7`,
	18,
).SetBuiltinToken(
	types2.TokenERC20,
	`(PoS) Dai Stablecoin`,
	`DAI`,
	`0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063`,
	18,
)
