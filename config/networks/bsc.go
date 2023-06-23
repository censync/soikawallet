//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var BSC = types.NewNetwork(
	types.BSC,
	`Binance Smart Chain`,
	`BNB`,
	18,
	true,
	&types.EVMConfig{
		ChainId:  0x38,
		DataFeed: "0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee",
	},
).SetDefaultRPC(
	`https://bsc-dataseed1.binance.org/`,
	`https://bscscan.com/`, // /block/ /address/ /tx/
).SetDataFeeds(map[types.CurrencyPair]string{
	{"BNB", "USD"}:   "0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee",
	{"BTC", "USD"}:   "0x264990fbd0a4796a3e3d8e37c4d5f87a3aca5ebf",
	{"ETH", "USD"}:   "0x9ef1b8c0e4f7dc8bf5719ea496883dc6401d5b2e",
	{"MATIC", "USD"}: "0x7ca57b0ca6367191c94c8914d7df09a57655905f",
},
).SetBuiltinToken(
	types.TokenBEP20,
	`Ethereum`,
	`ETH`,
	`0x2170ed0880ac9a755fd29b2688956bd959f933f8`,
	18,
).SetBuiltinToken(
	types.TokenBEP20,
	`BUSD`,
	`BUSD`,
	`0x55d398326f99059ff775485246999027b3197955`,
	18,
).SetBuiltinToken(
	types.TokenBEP20,
	`USD Coin`,
	`USDC`,
	`0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d`,
	18,
)
