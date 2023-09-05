//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Polygon = types.NewNetwork(
	mhda.MATIC,
	`Polygon (Testnet)`,
	`MATIC`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x13881,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon_testnet`,
	`https://mumbai.polygonscan.com/`, // /block/ /address/ /tx/
)
