package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

func (s *Wallet) RPC(dto *dto.GetRPCListByIndexDTO) *types.RPC {
	return s.getRPCProvider(types.NetworkType(dto.NetworkType)).RPC(dto.Index)
}

func (s *Wallet) AllRPC(dto *dto.GetRPCListByNetworkDTO) map[uint32]*types.RPC {
	return s.getRPCProvider(types.NetworkType(dto.NetworkType)).AllRPC()
}

func (s *Wallet) AddRPC(dto *dto.AddRPCDTO) error {
	provider := s.getRPCProvider(types.NetworkType(dto.NetworkType))

	index, err := provider.AddRPC(dto.Title, dto.Endpoint)

	if err != nil {
		return err
	}
	rpcConfig := provider.RPC(index)

	// TODO: Add atomic
	s.meta.AddRPCNode(types.NodeIndex{
		NetworkType: types.NetworkType(dto.NetworkType),
		Index:       index,
	}, rpcConfig)

	return nil
}

// not used
func (s *Wallet) RemoveRPC(dto *dto.RemoveRPCDTO) error {
	nodeIndex := types.NodeIndex{
		NetworkType: types.NetworkType(dto.NetworkType),
		Index:       dto.Index,
	}

	if len(s.meta.GetRPCAccountLinks(nodeIndex)) > 0 {
		return errors.New("cannot remove rpc, until linked accounts exists")
	}

	provider := s.getRPCProvider(types.NetworkType(dto.NetworkType))

	err := provider.RemoveRPC(dto.Index)

	if err != nil {
		return err
	}

	s.meta.RemoveRPCNode(nodeIndex)

	return nil
}

func (s *Wallet) AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error {
	nodeIndex := types.NodeIndex{
		NetworkType: types.NetworkType(dto.NetworkType),
		Index:       dto.NodeIndex,
	}
	return s.setAccountLinkRPC(nodeIndex, types.AccountIndex(dto.AccountIndex))
}

func (s *Wallet) setAccountLinkRPC(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) error {
	if s.getRPCProvider(nodeIndex.NetworkType).RPC(nodeIndex.Index) == nil {
		return errors.New("undefined node index")
	}

	err := s.meta.SetRPCAccountLink(nodeIndex, accountIndex)

	if err != nil {
		return err
	}

	// set for addresses
	for index, addr := range s.addresses {
		if addr.Network() == nodeIndex.NetworkType && addr.Account() == accountIndex {
			s.addresses[index].nodeIndex = nodeIndex.Index
		}
	}
	return nil
}

func (s *Wallet) RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error {
	return s.removeAccountLinkRPC(types.NetworkType(dto.NetworkType), types.AccountIndex(dto.AccountIndex))
}

func (s *Wallet) removeAccountLinkRPC(networkType types.NetworkType, accountIndex types.AccountIndex) error {
	var (
		nodeKey types.NodeIndex
		isExist bool
	)

	nodeInstance := s.getRPCProvider(nodeKey.NetworkType).RPC(nodeKey.Index)

	for _, addr := range s.addresses {
		if addr.Network() == networkType &&
			addr.Account() == accountIndex {
			nodeKey = types.NodeIndex{
				NetworkType: networkType,
				Index:       addr.nodeIndex,
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

	//s.removeAccountLinkRPC(networkType, accountIndex)

	s.meta.RemoveRPCAccountLink(nodeKey, accountIndex)

	return nil
}

func (s *Wallet) GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int {
	return s.meta.GetRPCAccountLinksCount(types.NetworkType(dto.NetworkType), dto.NodeIndex)
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
	ctx := types.NewRPCContext(types.NetworkType(dto.NetworkType), dto.NodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	return provider.GetRPCInfo(ctx)
}

func (s *Wallet) GetBaseCurrency(dto *dto.GetTokensByNetworkDTO) (*resp.BaseCurrency, error) {
	if !types.IsNetworkExists(types.NetworkType(dto.NetworkType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	rpcProvider := s.getRPCProvider(types.NetworkType(dto.NetworkType))

	if rpcProvider == nil {
		return nil, errors.New("cannot get RPC instance")
	}

	return &resp.BaseCurrency{
		Symbol:   rpcProvider.Currency(),
		Decimals: rpcProvider.Decimals(),
	}, nil
}

func (s *Wallet) GetAllTokensByNetwork(dto *dto.GetTokensByNetworkDTO) (*resp.AddressTokensListResponse, error) {
	if !types.IsNetworkExists(types.NetworkType(dto.NetworkType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	result := resp.AddressTokensListResponse{}

	rpcProvider := s.getRPCProvider(types.NetworkType(dto.NetworkType))

	if rpcProvider == nil {
		return nil, errors.New("cannot get RPC instance")
	}

	rpcTokens := rpcProvider.GetAllTokens()

	rpcBaseToken := rpcProvider.GetBaseToken()

	result[rpcBaseToken.Contract()] = &resp.AddressTokenEntry{
		Standard: uint8(rpcBaseToken.Standard()),
		Name:     rpcBaseToken.Name(),
		Symbol:   rpcBaseToken.Symbol(),
		Contract: rpcBaseToken.Contract(),
	}

	for contract, token := range rpcTokens {
		result[contract] = &resp.AddressTokenEntry{
			Standard: uint8(token.Standard()),
			Name:     token.Name(),
			Symbol:   token.Symbol(),
			Contract: token.Contract(),
		}
	}
	return &result, nil
}

func (s *Wallet) GetTokensByPath(dto *dto.GetAddressTokensByPathDTO) (*resp.AddressTokensListResponse, error) {
	addrPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addressLinkedTokenContracts, err := s.meta.GetAddressTokens(addrPath.Network(), addrPath.Account(), addrPath.AddressIndex())

	result := resp.AddressTokensListResponse{}

	rpcProvider := s.getRPCProvider(addrPath.Network())
	rpcBaseToken := rpcProvider.GetBaseToken()

	result[rpcBaseToken.Contract()] = &resp.AddressTokenEntry{
		Standard: uint8(rpcBaseToken.Standard()),
		Name:     rpcBaseToken.Name(),
		Symbol:   rpcBaseToken.Symbol(),
		Contract: rpcBaseToken.Contract(),
	}

	for _, tokenConfig := range addressLinkedTokenContracts {

		result[tokenConfig.Contract()] = &resp.AddressTokenEntry{
			Standard: uint8(tokenConfig.Standard()),
			Name:     tokenConfig.Name(),
			Symbol:   tokenConfig.Symbol(),
			Contract: tokenConfig.Contract(),
		}
	}

	return &result, nil
}

func (s *Wallet) GetToken(dto *dto.GetTokenDTO) (*resp.TokenConfig, error) {
	var (
		tokenConfig *types.TokenConfig
		err         error
	)

	if !types.IsNetworkExists(types.NetworkType(dto.NetworkType)) {
		return nil, errors.New("network is not exists in SLIP-44 list")
	}

	defaultNodeIndex := s.getRPCProvider(types.NetworkType(dto.NetworkType)).DefaultNodeId()

	ctx := types.NewRPCContext(types.NetworkType(dto.NetworkType), defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	tokenConfig = provider.GetTokenConfig(dto.Contract)
	if tokenConfig == nil {
		tokenConfig, err = provider.GetToken(ctx, dto.Contract)

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

	if !types.IsNetworkExists(types.NetworkType(dto.NetworkType)) {
		return errors.New("network is not exists in SLIP-44 list")
	}

	defaultNodeIndex := s.getRPCProvider(types.NetworkType(dto.NetworkType)).DefaultNodeId()

	ctx := types.NewRPCContext(types.NetworkType(dto.NetworkType), defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return err
	}

	if provider.IsTokenConfigExists(dto.Contract) {
		tokenConfig = provider.GetTokenConfig(dto.Contract)

		tokenIndex = types.TokenIndex{
			NetworkType: types.NetworkType(dto.NetworkType),
			Contract:    tokenConfig.Contract(),
		}
	} else {
		tokenConfig, err = provider.GetToken(ctx, dto.Contract)

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
			NetworkType: types.NetworkType(dto.NetworkType),
			Contract:    tokenConfig.Contract(),
		}
	}

	s.meta.AddTokenConfig(types.NetworkType(dto.NetworkType), tokenConfig)

	if dto.DerivationPath != "" {
		var addrPath *types.DerivationPath
		addrPath, err = types.ParsePath(dto.DerivationPath)
		if err != nil {
			return err
		}

		err = s.meta.SetTokenConfigAddressLink(tokenIndex, addrPath.Account(), addrPath.AddressIndex())
	}

	return err
}
