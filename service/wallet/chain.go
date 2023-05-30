package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/network"
)

func (s *Wallet) GetAllChains(dto *dto.GetChainsDTO) []*resp.ChainInfo {
	var result []*resp.ChainInfo
	for coiType, provider := range network.GetAll() {
		if provider.IsW3() {
			chainId := uint32(0)
			if provider.EVMConfig() != nil {
				chainId = provider.EVMConfig().ChainId
			}
			result = append(result, &resp.ChainInfo{
				CoinType: coiType,
				Name:     provider.Name(),
				ChainId:  chainId,
			})
		}
	}
	return result
}
