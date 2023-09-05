//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Moonbeam = types.NewNetwork(
	mhda.ETH,
	`Moonbeam`,
	`GLMR`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x504,
	},
).SetDefaultRPC(
	`https://1rpc.io/glmr`,
	`https://moonbeam.moonscan.io/`, // /block/ /address/ /tx/
)
