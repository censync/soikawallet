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
	"github.com/censync/soikawallet/service/core/internal/types"
)

var Tron = types.NewNetwork(
	mhda.TRX,
	`Tron`,
	`TRX`,
	6,
	10e6,
	"SUN",
	false,
	nil,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/tron`, // `http://3.225.171.164`,
	`https://tronscan.org/`,                 // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenTRC20,
	`Wrapped TRX`,
	`WTRX`,
	`TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR`,
	6,
).SetBuiltinToken(
	types.TokenTRC20,
	`Wrapped BitTorrent`,
	`WBTT`,
	`TKfjV9RNKJJCqPvBtK8L7Knykh7DNWvnYt`,
	6,
).SetBuiltinToken(
	types.TokenTRC20,
	`Bitcoin`,
	`BTC`,
	`TN3W4H6rK2ce4vX9YnFQHwKENnHjoxb3m9`,
	8,
).SetBuiltinToken(
	types.TokenTRC20,
	`Ethereum`,
	`ETH`,
	`THb4CqiFdwNHsWsQCs4JhzwjMWys4aqCbF`,
	18,
).SetBuiltinToken(
	types.TokenTRC20,
	`Tether USD`,
	`USDT`,
	`TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t`,
	6,
).SetBuiltinToken(
	types.TokenTRC20,
	`USD network`,
	`USDC`,
	`TEkxiTehnzSmSe2XqrBj4w32RUN966rdz8`,
	6,
).SetBuiltinToken(
	types.TokenTRC20,
	`Decentralized USD`,
	`USDD`,
	`TPYmHEhy5n8TCEfYGqW2rPxsghSfzghPDn`,
	18,
)
