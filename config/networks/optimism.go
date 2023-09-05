//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Optimism = types.NewNetwork(
	mhda.ETH,
	`Optimism`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0xa,
	},
).SetDefaultRPC(
	`https://mainnet.optimism.io`,
	`https://explorer.optimism.io`, // /block/ /address/ /tx/
)
