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

package core

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config/datafeed"
	"github.com/censync/soikawallet/types"
	"math"
)

func (s *Wallet) UpdateFiatCurrencies() map[string]float64 {
	loadedPairs := map[string]float64{}
	fiatPairs := datafeed.GetFiatDataFeeds(fiatTitle)
	if len(fiatPairs) > 0 {
		for index := range fiatPairs {
			pair := fiatPairs[index]
			if _, ok := loadedPairs[pair.Pair]; !ok {
				defaultNodeIndex := s.getRPCProvider(pair.ChainKey).DefaultNodeId()

				ctx := types.NewRPCContext(pair.ChainKey, defaultNodeIndex)

				provider, err := s.getNetworkProvider(ctx)

				if err == nil {
					value, decimals, err := provider.ChainLinkGetPrice(ctx, pair.Address)
					if err == nil {
						calculatedPrice := float64(value) / (math.Pow(10, float64(decimals)))
						s.currenciesFiat.Set(pair.Symbol, calculatedPrice, pair.Type, pair.Address, pair.ChainKey)
						loadedPairs[pair.Pair] = calculatedPrice
					}
				}
			}
		}
	}
	return loadedPairs
}

func (s *Wallet) GetFiatCurrency(dto *dto.GetFiatCurrencyDTO) (float64, string, string) {
	provider := s.getRPCProvider(dto.ChainKey)

	if currency := s.currenciesFiat.Get(provider.Currency()); currency != nil {
		return currency.Value(), s.currenciesFiat.Fiat(), s.currenciesFiat.Symbol()
	}
	return 0, s.currenciesFiat.Fiat(), s.currenciesFiat.Symbol()
}
