package currencies

import (
	"fmt"
	"github.com/censync/soikawallet/types"
	"time"
)

type FiatPair struct {
	sourceType DataFeedType
	source     string

	value float64

	network types.NetworkType
}

type FiatCurrencies struct {
	pairs     map[string]*FiatPair
	fiat      string
	updatedAt int64
}

func NewFiatCurrencies(fiat string) *FiatCurrencies {
	return &FiatCurrencies{pairs: map[string]*FiatPair{}, fiat: fiat}
}

func (f *FiatCurrencies) Set(symbol string, value float64, sourceType DataFeedType, source string, network types.NetworkType) {
	f.pairs[symbol] = &FiatPair{
		sourceType: sourceType,
		source:     source,
		value:      value,
		network:    network,
	}
	f.updatedAt = time.Now().Unix()
}

func (f *FiatCurrencies) Get(symbol string) *FiatPair {
	return f.pairs[symbol]
}

func (f *FiatCurrencies) AllSimple() map[string]string {
	result := map[string]string{}

	for title, pair := range f.pairs {
		result[title] = fmt.Sprintf("%.2f %s", pair.value, f.fiat)
	}
	return result
}

func (f *FiatCurrencies) Remove(symbol string) {
	if _, ok := f.pairs[symbol]; ok {
		delete(f.pairs, symbol)
	}
}
