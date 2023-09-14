package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/network"
)

func (s *Wallet) GetAllEvmChains(dto *dto.GetChainsDTO) []*resp.ChainInfo {
	var result []*resp.ChainInfo
	for chainKey, provider := range network.GetAll() {
		if provider.IsW3() {
			result = append(result, &resp.ChainInfo{
				ChainKey: chainKey,
				Name:     provider.Name(),
			})
		}
	}
	return result
}
