package core

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/types"
)

// GetTxReceipt returns transaction info by hash
func (s *Wallet) GetTxReceipt(dto *dto.GetTxReceiptDTO) (map[string]interface{}, error) {
	ctx := types.NewRPCContext(dto.ChainKey, dto.NodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	return provider.TxGetReceipt(ctx, dto.Hash)
}
