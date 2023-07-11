//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var Polygon = types.NewNetwork(
	types.Polygon,
	`Polygon`,
	`MATIC`,
	18,
	true,
	&types.EVMConfig{
		ChainId: 0x89,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon`, // https://polygon-rpc.com
	`https://polygonscan.com/`,                 // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenERC20,
	`(PoS) Tether USD`,
	`USDT`,
	`0xc2132d05d31c914a87c6611c10748aeb04b58e8f`,
	6,
).SetBuiltinToken(
	types.TokenERC20,
	`USD Network (PoS)`,
	`USDC`,
	`0x2791bca1f2de4661ed88a30c99a7a9449aa84174`,
	6,
).SetBuiltinToken(
	types.TokenERC20,
	`(PoS) Binance USD`,
	`BUSD`,
	`0xdAb529f40E671A1D4bF91361c21bf9f0C9712ab7`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`(PoS) Dai Stablecoin`,
	`DAI`,
	`0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063`,
	18,
)
