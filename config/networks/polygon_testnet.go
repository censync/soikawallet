//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var Polygon = types.NewNetwork(
	types.Polygon,
	`Polygon (Testnet)`,
	`MATIC`,
	18,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon_testnet`,
	`https://mumbai.polygonscan.com/`, // /block/ /address/ /tx/
)
