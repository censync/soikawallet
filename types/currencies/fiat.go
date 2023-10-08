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

package currencies

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"time"
)

type FiatPair struct {
	sourceType DataFeedType
	source     string

	value float64

	chainKey mhda.ChainKey
}

func (p *FiatPair) Value() float64 {
	return p.value
}

type FiatCurrencies struct {
	pairs     map[string]*FiatPair
	fiat      string
	symbol    string
	updatedAt int64
}

func (f *FiatCurrencies) Fiat() string {
	return f.fiat
}

func (f *FiatCurrencies) Symbol() string {
	return f.symbol
}

func NewFiatCurrencies(fiat, symbol string) *FiatCurrencies {
	return &FiatCurrencies{pairs: map[string]*FiatPair{}, fiat: fiat, symbol: symbol}
}

func (f *FiatCurrencies) Set(symbol string, value float64, sourceType DataFeedType, source string, chainKey mhda.ChainKey) {
	f.pairs[symbol] = &FiatPair{
		sourceType: sourceType,
		source:     source,
		value:      value,
		chainKey:   chainKey,
	}
	f.updatedAt = time.Now().Unix()
}

func (f *FiatCurrencies) Get(symbol string) *FiatPair {
	return f.pairs[symbol]
}

func (f *FiatCurrencies) Exists(symbol string) bool {
	_, exists := f.pairs[symbol]
	return exists
}

func (f *FiatCurrencies) AllSimple() map[string]string {
	result := map[string]string{}

	for title, pair := range f.pairs {
		result[title] = fmt.Sprintf("%.2f %s", pair.value, f.symbol)
	}
	return result
}

func (f *FiatCurrencies) Remove(symbol string) {
	if _, ok := f.pairs[symbol]; ok {
		delete(f.pairs, symbol)
	}
}
