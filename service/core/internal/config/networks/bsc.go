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

var BSC = types2.NewNetwork(
	mhda.BSC,
	`Binance Smart Chain`,
	`BNB`,
	18,
	1e9,
	"gwei",
	true,
	&types2.EVMConfig{
		ChainId: 0x38,
	},
).SetDefaultRPC(
	`https://bsc-dataseed1.binance.org/`,
	`https://bscscan.com/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types2.TokenBEP20,
	`Ethereum`,
	`ETH`,
	`0x2170ed0880ac9a755fd29b2688956bd959f933f8`,
	18,
).SetBuiltinToken(
	types2.TokenBEP20,
	`BUSD`,
	`BUSD`,
	`0x55d398326f99059ff775485246999027b3197955`,
	18,
).SetBuiltinToken(
	types2.TokenBEP20,
	`USD network`,
	`USDC`,
	`0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d`,
	18,
)
