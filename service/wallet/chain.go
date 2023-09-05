package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/network"
	"strconv"
)

func (s *Wallet) GetAllEvmChains(dto *dto.GetChainsDTO) []*resp.ChainInfo {
	var result []*resp.ChainInfo
	for chainKey, provider := range network.GetAll() {
		if provider.IsW3() {
			chainId := uint32(0)
			if provider.EVMConfig() != nil {
				chainId = provider.EVMConfig().ChainId
			}
			result = append(result, &resp.ChainInfo{
				ChainKey: chainKey,
				Name:     provider.Name(),
				ChainId:  "0x" + strconv.FormatUint(uint64(chainId), 16), // TODO: Update with mhda
			})
		}
	}
	return result
}
