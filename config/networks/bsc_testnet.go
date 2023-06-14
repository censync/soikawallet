//go:build testnet

package networks

import "github.com/censync/soikawallet/types"

var BSC = types.NewNetwork(
	types.BSC,
	`BSC (Testnet)`,
	`BNBBT`,
	18,
	true,
	&types.EVMConfig{
		ChainId:  0x61,
		DataFeed: "",
	},
).SetDefaultRPC(
	`https://data-seed-prebsc-1-s1.binance.org:8545/`,
	`https://testnet.bscscan.com/`, // /block/ /address/ /tx/
)
