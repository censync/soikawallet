//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var Tron = types.NewNetwork(
	types.Tron,
	`Tron`,
	`TRX`,
	6,
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
	`USD Coin`,
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
