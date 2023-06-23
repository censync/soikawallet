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
		ChainId:  0x89,
		DataFeed: "0xab594600376ec9fd91f8e885dadf0ce036862de0",
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/polygon`, // https://polygon-rpc.com
	`https://polygonscan.com/`,                 // /block/ /address/ /tx/
).SetDataFeeds(map[types.CurrencyPair]string{
	{"MATIC", "USD"}: "0xab594600376ec9fd91f8e885dadf0ce036862de0",
	{"BTC", "USD"}:   "0xc907e116054ad103354f2d350fd2514433d57f6f",
	{"ETH", "USD"}:   "0xf9680d99d6c9589e2a93a78a04a279e509205945",
	{"BNB", "USD"}:   "0x82a6c4af830caa6c97bb504425f6a66165c2c26e",
},
).SetBuiltinToken(
	types.TokenERC20,
	`(PoS) Tether USD`,
	`USDT`,
	`0xc2132d05d31c914a87c6611c10748aeb04b58e8f`,
	6,
).SetBuiltinToken(
	types.TokenERC20,
	`USD Coin (PoS)`,
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
