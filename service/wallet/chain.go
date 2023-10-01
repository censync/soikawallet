package wallet

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/wallet/internal/network"
)

func (s *Wallet) GetAllEvmW3Chains() []*resp.ChainInfo {
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
