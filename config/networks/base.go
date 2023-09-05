//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Base = types.NewNetwork(
	mhda.ETH,
	`Base`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x2105,
	},
).SetDefaultRPC(
	`https://mainnet.base.org`,
	`https://basescan.org/`, // /block/ /address/ /tx/
)
