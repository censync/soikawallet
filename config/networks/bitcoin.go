//go:build !testnet

package networks

import (
	"github.com/censync/soikawallet/types"
)

var Bitcoin = types.NewNetwork(
	types.Bitcoin,
	`Bitcoin`,
	`BTC`,
	8,
	1e8,
	"satoshi",
	false,
	nil,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://etherscan.io/`, // /block/ /address/ /tx/
)
