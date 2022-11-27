//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var Ethereum = types.NewNetwork(
	types.Ethereum,
	`Ethereum`,
	`ETH`,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://etherscan.io/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	`USDT`,
	`0xdAC17F958D2ee523a2206206994597C13D831ec7`,
).SetBuiltinToken(
	`USDC`,
	`0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48`,
).SetBuiltinToken(
	`BUSD`,
	`0x4Fabb145d64652a948d72533023f6E7A623C7C53`,
).SetBuiltinToken(
	`MATIC`,
	`0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0`,
).SetBuiltinToken(
	`DAI`,
	`0x6B175474E89094C44Da98b954EedeAC495271d0F`,
).SetBuiltinToken(
	`LINK`,
	`0x514910771AF9Ca656af840dff83E8264EcF986CA`,
).SetBuiltinToken(
	`SHIB`,
	`0x95aD61b0a150d79219dCF64E1E6Cc01f0B64C4cE`,
)
