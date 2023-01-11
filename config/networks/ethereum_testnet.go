//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var Ethereum = types.NewNetwork(
	types.Ethereum,
	`Ethereum (Testnet)`,
	`SepETH`,
	18,
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth_testnet2`,
	`https://sepolia.etherscan.io/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenERC20,
	`KOKToken`,
	`KOK_STG`,
	`0x91B333A8485737f9B93327483030f48526FaDc22`,
	18,
)
