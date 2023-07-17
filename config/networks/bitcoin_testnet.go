//go:build testnet

package networks

import (
	"github.com/censync/soikawallet/types"
)

var Bitcoin = types.NewNetwork(
	types.Bitcoin,
	`Bitcoin`,
	`BTC`,
	8,
	10e8,
	"satoshi",
	false,
	nil,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://www.blockchain.com/explorer/blocks/btc/`, // /block/ /address/ /tx/
)
