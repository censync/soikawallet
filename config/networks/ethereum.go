//go:build !testnet

package networks

import (
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/gas"
)

var Ethereum = types.NewNetwork(
	types.Ethereum,
	`Ethereum`,
	`ETH`,
	18,
	true,
	&types.EVMConfig{
		ChainId:  0x1,
		DataFeed: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",
	},
).SetGasCalculator(&gas.CalcEVML1V1{
	CalcOpts: &gas.CalcOpts{
		GasSuffix:     "gwei",
		TokenCurrency: 10e9,
		TokenSuffix:   "Eth",
	},
	Units:       21000,
	BaseFee:     0,
	PriorityFee: 0,
	MaxFee:      30e6,
},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://etherscan.io/`, // /block/ /address/ /tx/
).SetDataFeeds(map[types.CurrencyPair]string{
	{"ETH", "USD"}:   "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419",
	{"BTC", "USD"}:   "0xf4030086522a5beea4988f8ca5b36dbc97bee88c",
	{"BNB", "USD"}:   "0x14e613ac84a31f709eadbdf89c6cc390fdc9540a",
	{"MATIC", "USD"}: "0x7bac85a8a13a4bcd8abb3eb7d6b4d632c5a57676",
},
).SetBuiltinToken(
	types.TokenERC20,
	`Tether USD`,
	`USDT`,
	`0xdAC17F958D2ee523a2206206994597C13D831ec7`,
	6,
).SetBuiltinToken(
	types.TokenERC20,
	`USD Coin`,
	`USDC`,
	`0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`,
	6,
).SetBuiltinToken(
	types.TokenERC20,
	`Binance USD`,
	`BUSD`,
	`0x4Fabb145d64652a948d72533023f6E7A623C7C53`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`Matic Token`,
	`MATIC`,
	`0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`Dai Stablecoin`,
	`DAI`,
	`0x6B175474E89094C44Da98b954EedeAC495271d0F`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`ChainLink Token`,
	`LINK`,
	`0x514910771AF9Ca656af840dff83E8264EcF986CA`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`SHIBA INU`,
	`SHIB`,
	`0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE`,
	18,
)
