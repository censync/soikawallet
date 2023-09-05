//go:build testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
)

var BSC = types.NewNetwork(
	mhda.BSC,
	`BSC (Testnet)`,
	`BNBBT`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x61,
	},
).SetDefaultRPC(
	`https://data-seed-prebsc-1-s1.binance.org:8545/`,
	`https://testnet.bscscan.com/`, // /block/ /address/ /tx/
)
