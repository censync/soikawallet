//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var Ethereum = types.NewNetwork(
	types.Ethereum,
	`Ethereum (Testnet)`,
	`SepETH`,
	18,
	true,
	&types.EVMConfig{
		ChainId: 11155111,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth_testnet2`,
	`https://sepolia.etherscan.io/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenERC20,
	`CenTest Token v1`,
	`CEN_TV1`,
	`0x91B268bd44c6a16b2E518060b44eFF33cB17f84d`,
	18,
)
