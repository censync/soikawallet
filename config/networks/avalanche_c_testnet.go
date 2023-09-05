//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

// https://support.avax.network/en/articles/7004986-what-derivation-paths-does-avalanche-use
var AvalancheC = types.NewNetwork(
	mhda.ETH,
	`Avalanche C-Chain (Testnet)`,
	`AVAX`,
	18,
	1e9,
	"nAVAX",
	true,
	&types.EVMConfig{
		ChainId: 0xa869,
	},
).SetDefaultRPC(
	`https://api.avax-test.network/ext/bc/C/rpc`,
	`https://testnet.snowtrace.io/`, // /block/ /address/ /tx/
)
