package currencies

import "github.com/censync/soikawallet/types"

type TokenPair struct {
	SymbolSrc  string
	SymbolDst  string
	AddressSrc string
	AddressDst string
	// source datafeed source
	DataFeedType uint8
	DataFeed     string
	Network      types.NetworkType
	Standard     types.TokenStandard
}

type TokenCurrencies struct {
	pairs     map[TokenPair]float64
	updatedAt uint64
}

func (c *TokenCurrencies) Set() {

}
