//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Tron = types.NewNetwork(
	mhda.TRX,
	`Tron (Testnet)`,
	`TRX`,
	6,
	10e6,
	"SUN",
	false,
	nil,
).SetDefaultRPC(
	`https://api.shasta.trongrid.io`, // https://rpc.soikawallet.app:8431/tron_testnet
	`https://shasta.tronscan.org/`,   // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenTRC20,
	`USDT on Shasta Test Net`,
	`USDT`,
	`TG3XXyExBkPp9nzdajDZsozEu4BkaSJozs`,
	6,
)
