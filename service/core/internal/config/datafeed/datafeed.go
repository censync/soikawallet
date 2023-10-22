// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package datafeed

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types/currencies"
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
