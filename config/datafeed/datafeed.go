package datafeed

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types/currencies"
	"strings"
)

type FeedSource struct {
	Type     currencies.DataFeedType
	ChainKey mhda.ChainKey
	Symbol   string
	Pair     string
	Address  string
}

func GetFiatDataFeeds(fiat string) []FeedSource {
	var result []FeedSource
	for dataFeedType, dataFeedTypes := range evmFiat {
		for chainKey, networkTypes := range dataFeedTypes {
			for pair, address := range networkTypes {
				symbol := ``
				if strings.HasPrefix(pair, fiat+`_`) {
					symbol = strings.Trim(pair, fiat+`_`)
				} else if strings.HasSuffix(pair, `_`+fiat) {
					symbol = strings.Trim(pair, `_`+fiat)
				}
				result = append(result, FeedSource{
					Type:     dataFeedType,
					ChainKey: chainKey,
					Symbol:   symbol,
					Pair:     pair,
					Address:  address,
				})
			}
		}
	}
	return result
}
