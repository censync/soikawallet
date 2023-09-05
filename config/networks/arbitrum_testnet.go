//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var ArbitrumOne = types.NewNetwork(
	mhda.ETH,
	`Arbitrum Sepolia`,
	`AGOR`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x66eee,
	},
).SetDefaultRPC(
	`https://sepolia-rollup.arbitrum.io/rpc`,
	`https://sepolia.arbiscan.io`, // /block/ /address/ /tx/
)
