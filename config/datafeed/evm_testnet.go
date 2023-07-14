//go:build testnet

package datafeed

import (
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/currencies"
)

var evmFiat = map[currencies.DataFeedType]map[types.NetworkType]map[string]string{
	currencies.FeedChainLink: {
		types.Ethereum: {
			"ETH_USD": "0x694AA1769357215DE4FAC081bf1f309aDC325306",
			"BTC_USD": "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43",
		},
		types.Polygon: {
			"MATIC_USD": "0xd0D5e3DB44DE05E9F294BB0a3bEEaF030DE24Ada",
			"BTC_USD":   "0x007A22900a3B98143368Bd5906f8E17e9867581b",
			"ETH_USD":   "0x0715A7794a1dc8e42615F059dD6e406A6594651A",
		},
		types.BSC: {
			"BNB_USD":   "0x2514895c72f50D8bd4B4F9b1110F0D6bD2c97526",
			"BTC_USD":   "0x5741306c21795FdCBb9b265Ea0255F499DFe515C",
			"ETH_USD":   "0x143db3CEEfbdfe5631aDD3E50f7614B6ba708BA7",
			"MATIC_USD": "0x957Eb0316f02ba4a9De3D308742eefd44a3c1719",
		},
	},
}
