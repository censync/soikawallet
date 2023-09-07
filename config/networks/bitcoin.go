//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Bitcoin = types.NewNetwork(
	mhda.BTC,
	`Bitcoin`,
	`BTC`,
	8,
	1e8,
	"satoshi",
	false,
	nil,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/btc`,
	`https://etherscan.io/`, // /block/ /address/ /tx/
)
