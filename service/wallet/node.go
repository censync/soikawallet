package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error {
	return s.setAccountLinkRPC(types.CoinType(dto.CoinType), types.AccountIndex(dto.AccountIndex), dto.NodeIndex)
}

func (s *Wallet) setAccountLinkRPC(coinType types.CoinType, accountIndex types.AccountIndex, nodeIndex uint32) error {
	if s.getRPCProvider(coinType).RPC(nodeIndex) == nil {
		return errors.New("undefined node index")
	}

	err := s.meta.SetRPCAccountLink(coinType, accountIndex, nodeIndex)

	if err != nil {
		return err
	}

	// set for addresses
	for index, address := range s.addresses {
		if address.CoinType() == coinType && address.Account() == accountIndex {
			s.addresses[index].nodeIndex = nodeIndex
		}
	}
	return nil
}

// RemoveRPCLinkedAccountDTO
func (s *Wallet) RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error {
	return s.removeAccountLinkRPC(types.CoinType(dto.CoinType), types.AccountIndex(dto.AccountIndex))
}

func (s *Wallet) removeAccountLinkRPC(coinType types.CoinType, accountIndex types.AccountIndex) error {
	var (
		nodeKey types.NodeIndex
		isExist bool
	)
	for _, address := range s.addresses {
		if address.CoinType() == coinType &&
			address.Account() == accountIndex {
			nodeKey = types.NodeIndex{
				CoinType: coinType,
				Index:    address.nodeIndex,
			}
			isExist = true
			break
		}
	}
	if !isExist {
		return errors.New("account is not found")
	}

	nodeInstance := s.getRPCProvider(nodeKey.CoinType).RPC(nodeKey.Index)

	if nodeInstance != nil && nodeInstance.IsDefault() {
		return errors.New("cannot unlink default node")
	}

	s.meta.RemoveRPCAccountLink(nodeKey, accountIndex)

	return s.setAccountLinkRPC(coinType, accountIndex, 0)
}

func (s *Wallet) GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int {
	return s.meta.GetRPCAccountLinksCount(types.CoinType(dto.CoinType), dto.NodeIndex)
}

func (s *Wallet) setAddressNode(path *types.DerivationPath, nodeIndex uint32) error {
	_, ok := s.addresses[path.String()]
	if !ok {
		return errors.New("addr is not found")
	}

	s.addresses[path.String()].nodeIndex = nodeIndex

	return nil
}

func (s *Wallet) GetRPCInfo(dto *dto.GetRPCInfoDTO) (map[string]interface{}, error) {
	ctx := types.NewRPCContext(types.CoinType(dto.CoinType), dto.NodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	return provider.GetRPCInfo(ctx)
}
