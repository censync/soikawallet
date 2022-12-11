//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var Ethereum = types.NewNetwork(
	types.Ethereum,
	`Ethereum (Testnet)`,
	`SepETH`,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth_testnet2`,
	`https://sepolia.etherscan.io/`, // /block/ /address/ /tx/
)
