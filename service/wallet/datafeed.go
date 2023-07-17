package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config/datafeed"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) UpdateFiatCurrencies() map[string]float64 {
	loadedPairs := map[string]float64{}
	fiatPairs := datafeed.GetFiatDataFeeds(fiatTitle)
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
						s.currenciesFiat.Set(pair.Symbol, value, pair.Type, pair.Address, pair.Network)
						loadedPairs[pair.Pair] = value
					}
				}
			}
		}
	}
	return loadedPairs
}

func (s *Wallet) GetFiatCurrency(dto *dto.GetFiatCurrencyDTO) (float64, string, string) {
	provider := s.getRPCProvider(types.NetworkType(dto.NetworkType))

	if currency := s.currenciesFiat.Get(provider.Currency()); currency != nil {
		return currency.Value(), s.currenciesFiat.Fiat(), s.currenciesFiat.Symbol()
	}
	return 0, s.currenciesFiat.Fiat(), s.currenciesFiat.Symbol()
}
