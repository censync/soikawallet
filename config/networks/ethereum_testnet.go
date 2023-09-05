//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Ethereum = types.NewNetwork(
	mhda.ETH,
	`Ethereum (Testnet)`,
	`ETH`, // SepETH
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0xaa36a7,
	},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth_testnet`,
	`https://sepolia.etherscan.io/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenERC20,
	`CenTest Claimable Token v4`,
	`CEN_TV4`,
	`0x8D2973D91C48540E9b7d1175885D97f38D03d0e8`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`CenTest Claimable Token v5`,
	`CEN_TV5`,
	`0x73F5Eb3092bd3D79D9b15EcEB1C560a72969142D`,
	18,
).SetBuiltinToken(
	types.TokenERC20,
	`CenTest Claimable Token v6`,
	`CEN_TV6`,
	`0xE639832e14644c273c1e51667947dFFB8B30EA6E`,
	18,
)
