//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var Tron = types.NewNetwork(
	types.Tron,
	`Tron`,
	`TRX`,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/tron`, // `http://3.225.171.164`,
	`https://tronscan.org/`,                 // /block/ /address/ /tx/
).SetBuiltinToken(
	`WTRX`,
	`TNUC9Qb1rRpS5CbWLmNMxXBjyFoydXjWFR`,
).SetBuiltinToken(
	`WBTT`,
	`TKfjV9RNKJJCqPvBtK8L7Knykh7DNWvnYt`,
).SetBuiltinToken(
	`BTC`,
	`TN3W4H6rK2ce4vX9YnFQHwKENnHjoxb3m9`,
).SetBuiltinToken(
	`ETH`,
	`THb4CqiFdwNHsWsQCs4JhzwjMWys4aqCbF`,
).SetBuiltinToken(
	`USDT`,
	`TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t`,
).SetBuiltinToken(
	`USDC`,
	`TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t`,
).SetBuiltinToken(
	`USDD`,
	`TPYmHEhy5n8TCEfYGqW2rPxsghSfzghPDn`,
)
