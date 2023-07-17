package datafeed

import (
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/currencies"
	"regexp"
	"strings"
)

var rxPairCheck = regexp.MustCompile("")

type FeedSource struct {
	Type    currencies.DataFeedType
	Network types.NetworkType
	Symbol  string
	Pair    string
	Address string
}

func GetFiatDataFeeds(fiat string) []FeedSource {
	var result []FeedSource
	for dataFeedType, dataFeedTypes := range evmFiat {
		for networkType, networkTypes := range dataFeedTypes {
			for pair, address := range networkTypes {
				symbol := ``
				if strings.HasPrefix(pair, fiat+`_`) {
					symbol = strings.Trim(pair, fiat+`_`)
				} else if strings.HasSuffix(pair, `_`+fiat) {
					symbol = strings.Trim(pair, `_`+fiat)
				}
				result = append(result, FeedSource{
					Type:    dataFeedType,
					Network: networkType,
					Symbol:  symbol,
					Pair:    pair,
					Address: address,
				})
			}
		}
	}
	return result
}
