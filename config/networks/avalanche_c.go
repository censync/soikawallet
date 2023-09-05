//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

// https://support.avax.network/en/articles/7004986-what-derivation-paths-does-avalanche-use
var AvalancheC = types.NewNetwork(
	mhda.ETH,
	`Avalanche C-Chain`,
	`AVAX`,
	18,
	1e9,
	"nAVAX",
	true,
	&types.EVMConfig{
		ChainId: 0xa86a,
	},
).SetDefaultRPC(
	`https://api.avax.network/ext/bc/C/rpc`,
	`https://snowtrace.io/`, // /block/ /address/ /tx/
)
