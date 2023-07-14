//go:build testnet

package networks

import (
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/gas"
)

var Bitcoin = types.NewNetwork(
	types.Bitcoin,
	`Bitcoin`,
	`BTC`,
	8,
	false,
	nil,
).SetGasCalculator(&gas.CalcBTCL1V1{
	CalcOpts: &gas.CalcOpts{
		GasSuffix:     "Satoshi",
		TokenCurrency: 10e8,
		TokenSuffix:   "BTC",
	},
},
).SetDefaultRPC(
	`https://rpc.soikawallet.app:8431/eth`,
	`https://www.blockchain.com/explorer/blocks/btc/`, // /block/ /address/ /tx/
)
