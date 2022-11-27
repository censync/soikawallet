package wallet

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) RPC(dto *dto.GetRPCListByIndexDTO) *types.RPC {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).RPC(dto.Index)
}

func (s *Wallet) AllRPC(dto *dto.GetRPCListByCoinDTO) map[uint32]*types.RPC {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).AllRPC()
}

func (s *Wallet) AddRPC(dto *dto.AddRPCDTO) (uint32, error) {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).AddRPC(dto.Title, dto.Endpoint)
}

func (s *Wallet) RemoveRPC(dto *dto.RemoveRPCDTO) error {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).RemoveRPC(dto.Index)
}
