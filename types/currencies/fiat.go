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
