package wallet

import (
	"github.com/censync/soikawallet/config/datafeed"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) UpdateFiatCurrencies() map[string]float64 {
	loadedPairs := map[string]float64{}
	fiatPairs := datafeed.GetFiatDataFeeds(fiatSymbol)
	if len(fiatPairs) > 0 {
		for index := range fiatPairs {
			pair := fiatPairs[index]
			if _, ok := loadedPairs[pair.Pair]; !ok {
				defaultNodeIndex := s.getRPCProvider(pair.Network).DefaultNodeId()

				ctx := types.NewRPCContext(pair.Network, defaultNodeIndex)

				provider, err := s.getNetworkProvider(ctx)

				if err == nil {
					value, err := provider.ChainLinkGetPrice(ctx, pair.Address)
					if err == nil {
						s.currenciesFiat.Set(pair.Pair, value, pair.Type, pair.Address, pair.Network)
						loadedPairs[pair.Pair] = value
					}
				}
			}
		}
	}
	return loadedPairs
}
