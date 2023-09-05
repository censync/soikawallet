//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var Moonbeam = types.NewNetwork(
	mhda.ETH,
	`Moonbeam testnet`,
	`DEV`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x507,
	},
).SetDefaultRPC(
	`https://rpc.api.moonbase.moonbeam.network`,
	`https://moonbase.moonscan.io/`, // /block/ /address/ /tx/
)
