//go:build !testnet

package networks

import "github.com/censync/soikawallet/types"

var BSC = types.NewNetwork(
	types.BSC,
	`Binance Smart Chain`,
	`BNB`,
	18,
	10e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x38,
	},
).SetDefaultRPC(
	`https://bsc-dataseed1.binance.org/`,
	`https://bscscan.com/`, // /block/ /address/ /tx/
).SetBuiltinToken(
	types.TokenBEP20,
	`Ethereum`,
	`ETH`,
	`0x2170ed0880ac9a755fd29b2688956bd959f933f8`,
	18,
).SetBuiltinToken(
	types.TokenBEP20,
	`BUSD`,
	`BUSD`,
	`0x55d398326f99059ff775485246999027b3197955`,
	18,
).SetBuiltinToken(
	types.TokenBEP20,
	`USD network`,
	`USDC`,
	`0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d`,
	18,
)
