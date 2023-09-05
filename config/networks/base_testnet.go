//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

// https://docs.base.org/network-information
var Base = types.NewNetwork(
	mhda.ETH,
	`Base Goerli`,
	`ETH`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x14a33,
	},
).SetDefaultRPC(
	`https://goerli.base.org`,
	`https://goerli.basescan.org`, // /block/ /address/ /tx/
)
