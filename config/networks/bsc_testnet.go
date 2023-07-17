//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var BSC = types.NewNetwork(
	types.BSC,
	`BSC (Testnet)`,
	`BNBBT`,
	18,
	10e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x61,
	},
).SetDefaultRPC(
	`https://data-seed-prebsc-1-s1.binance.org:8545/`,
	`https://testnet.bscscan.com/`, // /block/ /address/ /tx/
)
