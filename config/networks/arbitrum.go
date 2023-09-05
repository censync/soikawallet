//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var ArbitrumOne = types.NewNetwork(
	mhda.ETH,
	`Arbitrum One`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0xa4b1,
	},
).SetDefaultRPC(
	`https://arb1.arbitrum.io/rpc`,
	`https://arbiscan.io/`, // /block/ /address/ /tx/
)
