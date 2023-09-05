//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

// https://community.optimism.io/docs/useful-tools/networks/#parameters-for-node-operators-2
var Optimism = types.NewNetwork(
	mhda.ETH,
	`OP Sepolia`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0xaa37dc,
	},
).SetDefaultRPC(
	`https://sepolia.optimism.io`,
	``, // /block/ /address/ /tx/
)
