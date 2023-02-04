package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) RPC(dto *dto.GetRPCListByIndexDTO) *types.RPC {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).RPC(dto.Index)
}

func (s *Wallet) AllRPC(dto *dto.GetRPCListByCoinDTO) map[uint32]*types.RPC {
	return s.getRPCProvider(types.CoinType(dto.CoinType)).AllRPC()
}

func (s *Wallet) AddRPC(dto *dto.AddRPCDTO) error {
	provider := s.getRPCProvider(types.CoinType(dto.CoinType))

	index, err := provider.AddRPC(dto.Title, dto.Endpoint)

	if err != nil {
		return err
	}
	rpcConfig := provider.RPC(index)

	// TODO: Add atomic
	s.meta.AddRPCNode(types.NodeIndex{
		CoinType: types.CoinType(dto.CoinType),
		Index:    index,
	}, rpcConfig)

	return nil
}

// not used
func (s *Wallet) RemoveRPC(dto *dto.RemoveRPCDTO) error {
	nodeIndex := types.NodeIndex{
		CoinType: types.CoinType(dto.CoinType),
		Index:    dto.Index,
	}

	if len(s.meta.GetRPCAccountLinks(nodeIndex)) > 0 {
		return errors.New("cannot remove rpc, until linked accounts exists")
	}

	provider := s.getRPCProvider(types.CoinType(dto.CoinType))

	err := provider.RemoveRPC(dto.Index)

	if err != nil {
		return err
	}

	s.meta.RemoveRPCNode(nodeIndex)

	return nil
}

func (s *Wallet) AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error {
	nodeIndex := types.NodeIndex{
		CoinType: types.CoinType(dto.CoinType),
		Index:    dto.NodeIndex,
	}
	return s.setAccountLinkRPC(nodeIndex, types.AccountIndex(dto.AccountIndex))
}

func (s *Wallet) setAccountLinkRPC(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) error {
	if s.getRPCProvider(nodeIndex.CoinType).RPC(nodeIndex.Index) == nil {
		return errors.New("undefined node index")
	}

	err := s.meta.SetRPCAccountLink(nodeIndex, accountIndex)

	if err != nil {
		return err
	}

	// set for addresses
	for index, addr := range s.addresses {
		if addr.CoinType() == nodeIndex.CoinType && addr.Account() == accountIndex {
			s.addresses[index].nodeIndex = nodeIndex.Index
		}
	}
	return nil
}

func (s *Wallet) RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error {
	return s.removeAccountLinkRPC(types.CoinType(dto.CoinType), types.AccountIndex(dto.AccountIndex))
}

func (s *Wallet) removeAccountLinkRPC(coinType types.CoinType, accountIndex types.AccountIndex) error {
	var (
		nodeKey types.NodeIndex
		isExist bool
	)

	nodeInstance := s.getRPCProvider(nodeKey.CoinType).RPC(nodeKey.Index)

	for _, addr := range s.addresses {
		if addr.CoinType() == coinType &&
			addr.Account() == accountIndex {
			nodeKey = types.NodeIndex{
				CoinType: coinType,
				Index:    addr.nodeIndex,
			}
			isExist = true
			break
		}
	}

	if !isExist {
		return errors.New("account is not found")
	}

	if nodeInstance != nil && nodeInstance.IsDefault() {
		return errors.New("cannot unlink default node")
	}

	//s.removeAccountLinkRPC(coinType, accountIndex)

	s.meta.RemoveRPCAccountLink(nodeKey, accountIndex)

	return nil
}

func (s *Wallet) GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int {
	return s.meta.GetRPCAccountLinksCount(types.CoinType(dto.CoinType), dto.NodeIndex)
}

func (s *Wallet) setAddressNode(path *types.DerivationPath, nodeIndex uint32) error {
	_, ok := s.addresses[path.String()]
	if !ok {
		return errors.New("address is not found")
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

func (s *Wallet) GetBaseCurrency(dto *dto.GetTokensByNetworkDTO) (*resp.BaseCurrency, error) {
	if !types.IsCoinExists(types.CoinType(dto.CoinType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	rpcProvider := s.getRPCProvider(types.CoinType(dto.CoinType))

	if rpcProvider == nil {
		return nil, errors.New("cannot get RPC instance")
	}

	return &resp.BaseCurrency{
		Symbol:   rpcProvider.Currency(),
		Decimals: rpcProvider.Decimals(),
	}, nil
}

func (s *Wallet) GetTokensByNetwork(dto *dto.GetTokensByNetworkDTO) (*resp.AddressTokensListResponse, error) {
	if !types.IsCoinExists(types.CoinType(dto.CoinType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	result := resp.AddressTokensListResponse{}

	rpcProvider := s.getRPCProvider(types.CoinType(dto.CoinType))

	if rpcProvider == nil {
		return nil, errors.New("cannot get RPC instance")
	}

	rpcTokens := rpcProvider.AllTokens()

	for contract, token := range rpcTokens {
		result[contract] = &resp.AddressTokenEntry{
			Name:     token.Name(),
			Symbol:   token.Symbol(),
			Contract: token.Contract(),
		}
	}
	return &result, nil
}

func (s *Wallet) GetToken(dto *dto.GetTokenDTO) (*resp.TokenConfig, error) {
	var (
		tokenConfig *types.TokenConfig
		err         error
	)

	if !types.IsCoinExists(types.CoinType(dto.CoinType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	defaultNodeIndex := s.getRPCProvider(types.CoinType(dto.CoinType)).DefaultNodeId()

	ctx := types.NewRPCContext(types.CoinType(dto.CoinType), defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	tokenConfig = provider.GetTokenConfig(dto.Contract)
	if tokenConfig == nil {
		tokenConfig, err = provider.GetERC20Token(ctx, dto.Contract)

		if err != nil {
			return nil, err
		}

	}

	return &resp.TokenConfig{
		Standard: uint8(tokenConfig.Standard()),
		Name:     tokenConfig.Name(),
		Symbol:   tokenConfig.Symbol(),
		Contract: tokenConfig.Contract(),
		Decimals: tokenConfig.Decimals(),
	}, nil
}

func (s *Wallet) UpsertToken(dto *dto.AddTokenDTO) error {
	var (
		tokenConfig *types.TokenConfig
		tokenIndex  types.TokenIndex
		err         error
	)

	if !types.IsCoinExists(types.CoinType(dto.CoinType)) {
		return errors.New("network is not exists in SLIP-44 list")
	}

	defaultNodeIndex := s.getRPCProvider(types.CoinType(dto.CoinType)).DefaultNodeId()

	ctx := types.NewRPCContext(types.CoinType(dto.CoinType), defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return err
	}

	if provider.IsTokenConfigExists(dto.Contract) {
		tokenConfig = provider.GetTokenConfig(dto.Contract)

		tokenIndex = types.TokenIndex{
			CoinType: types.CoinType(dto.CoinType),
			Contract: tokenConfig.Contract(),
		}
	} else {
		tokenConfig, err = provider.GetERC20Token(ctx, dto.Contract)

		if err != nil {
			return err
		}

		tokenConfig, err = provider.AddTokenConfig(
			tokenConfig.Standard(),
			tokenConfig.Name(),
			tokenConfig.Symbol(),
			tokenConfig.Contract(),
			tokenConfig.Decimals(),
		)

		if err != nil {
			return err
		}

		tokenIndex = types.TokenIndex{
			CoinType: types.CoinType(dto.CoinType),
			Contract: tokenConfig.Contract(),
		}

		s.meta.AddTokenConfig(tokenIndex, tokenConfig)
	}

	if dto.DerivationPath != "" {
		addrPath, err := types.ParsePath(dto.DerivationPath)
		if err != nil {
			return err
		}

		err = s.meta.SetTokenConfigAddressLink(tokenIndex, addrPath.Account(), addrPath.AddressIndex())
	}

	return err
}
