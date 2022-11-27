//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var Polygon = types.NewNetwork(
	types.Polygon,
	`Polygon`,
	`MATIC`,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon`, // https://polygon-rpc.com
	`https://polygonscan.com/`,                 // /block/ /address/ /tx/
).SetBuiltinToken(
	`USDT`,
	`0xc2132d05d31c914a87c6611c10748aeb04b58e8f`,
).SetBuiltinToken(
	`USDC`,
	`0x2791bca1f2de4661ed88a30c99a7a9449aa84174`,
).SetBuiltinToken(
	`BUSD`,
	`0x9C9e5fD8bbc25984B178FdCE6117Defa39d2db39`,
).SetBuiltinToken(
	`DAI`,
	`0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063`,
)
