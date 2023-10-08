// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

var (
	errMetaRPCCannotRemoveLinkExists = errors.New("cannot remove rpc, until linked accounts exists")
	errRPCNotExists                  = errors.New("cannot get RPC instance")
	errAddrNotFound                  = errors.New("address not found")
	errChainKeyNotSupported          = errors.New("network is not supported")
)

func (s *Wallet) RPC(dto *dto.GetRPCListByIndexDTO) *types.RPC {
	return s.getRPCProvider(dto.ChainKey).RPC(dto.Index)
}

func (s *Wallet) AllRPC(dto *dto.GetRPCListByNetworkDTO) map[uint32]*types.RPC {
	return s.getRPCProvider(dto.ChainKey).AllRPC()
}

func (s *Wallet) AddRPC(dto *dto.AddRPCDTO) error {
	provider := s.getRPCProvider(dto.ChainKey)

	index, err := provider.AddRPC(dto.Title, dto.Endpoint)

	if err != nil {
		return err
	}
	rpcConfig := provider.RPC(index)

	// TODO: Add atomic
	return s.meta.AddRPCNode(types.NodeIndex{
		ChainKey: dto.ChainKey,
		Index:    index,
	}, rpcConfig)
}

// not used yet
func (s *Wallet) RemoveRPC(dto *dto.RemoveRPCDTO) error {
	nodeIndex := types.NodeIndex{
		ChainKey: dto.ChainKey,
		Index:    dto.Index,
	}

	if len(s.meta.GetRPCAccountLinks(nodeIndex)) > 0 {
		return errMetaRPCCannotRemoveLinkExists
	}

	provider := s.getRPCProvider(dto.ChainKey)

	err := provider.RemoveRPC(dto.Index)

	if err != nil {
		return err
	}

	return s.meta.RemoveRPCNode(nodeIndex)
}

func (s *Wallet) AccountLinkRPCSet(dto *dto.SetRPCLinkedAccountDTO) error {
	/*nodeIndex := types.NodeIndex{
		ChainKey: dto.ChainKey,
		Index:    dto.NodeIndex,
	}
	return s.setAccountLinkRPC(nodeIndex, types.AccountIndex(dto.AccountIndex))*/

	return nil
}

func (s *Wallet) setAccountLinkRPC(nodeIndex types.NodeIndex, accountIndex mhda.AccountIndex) error {
	/*if s.getRPCProvider(nodeIndex.ChainKey).RPC(nodeIndex.Index) == nil {
		return errors.New("undefined node index")
	}

	err := s.meta.SetRPCAddressLink(accountIndex, nodeIndex)

	if err != nil {
		return err
	}

	// set for addresses
	for index, addr := range s.meta.Addresses() {
		if addr.MHDA().Chain().key() == nodeIndex.ChainKey && addr.Account() == accountIndex {
			s.addresses[index].nodeIndex = nodeIndex.Index
		}
	}*/
	return nil
}

func (s *Wallet) RemoveAccountLinkRPC(dto *dto.RemoveRPCLinkedAccountDTO) error {
	//return s.removeAccountLinkRPC(dto.ChainKey, types.AccountIndex(dto.AccountIndex))
	return nil
}

func (s *Wallet) removeAccountLinkRPC(chainKey mhda.ChainKey, addrKey string) error {
	/* var (
		nodeKey types.NodeIndex
		isExist bool
	)

	nodeInstance := s.getRPCProvider(chainKey).RPC(nodeKey.Index)

	for _, addr := range s.meta.Addresses() {
		if addr.Network() == networkType &&
			addr.Account() == accountIndex {
			nodeKey = types.NodeIndex{
				ChainKey: addr.MHDA().Chain().key(),
				Index:    addr.NodeIndex(),
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

	s.meta.RemoveRPCAccountLink(accountIndex, nodeKey)
	*/
	return nil
}

func (s *Wallet) GetRPCLinkedAccountCount(dto *dto.GetRPCLinkedAccountCountDTO) int {
	return s.meta.GetRPCAccountLinksCount(types.NodeIndex{
		ChainKey: dto.ChainKey,
		Index:    dto.NodeIndex,
	})
}

func (s *Wallet) setAddressNode(addrKey string, nodeIndex uint32) error {
	addrPath, err := mhda.ParseNSS(addrKey)
	if err != nil {
		return err
	}

	addr := s.meta.GetAddress(addrPath.NSS())

	if addr == nil {
		return errors.New("address is not found")
	}

	//s.addresses[path.String()].nodeIndex = nodeIndex

	// s.meta.SetRPCAddressLink()

	return nil
}

func (s *Wallet) GetRPCInfo(dto *dto.GetRPCInfoDTO) (map[string]interface{}, error) {
	ctx := types.NewRPCContext(dto.ChainKey, dto.NodeIndex)
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return map[string]interface{}{}, err
	}
	return provider.GetRPCInfo(ctx)
}

func (s *Wallet) GetBaseCurrency(dto *dto.GetTokensByNetworkDTO) (*resp.BaseCurrency, error) {
	if !types.IsNetworkExists(dto.ChainKey) {
		return nil, errChainKeyNotSupported
	}

	rpcProvider := s.getRPCProvider(dto.ChainKey)

	if rpcProvider == nil {
		return nil, errRPCNotExists
	}

	return &resp.BaseCurrency{
		Symbol:   rpcProvider.Currency(),
		Decimals: rpcProvider.Decimals(),
	}, nil
}

func (s *Wallet) GetAllTokensByNetwork(dto *dto.GetTokensByNetworkDTO) (*resp.AddressTokensListResponse, error) {
	if !types.IsNetworkExists(dto.ChainKey) {
		return nil, errChainKeyNotSupported
	}

	result := resp.AddressTokensListResponse{}

	rpcProvider := s.getRPCProvider(dto.ChainKey)

	if rpcProvider == nil {
		return nil, errRPCNotExists
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
	addrPath, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return nil, err
	}

	addr := s.meta.GetAddress(addrPath.NSS())

	if addr == nil {
		return nil, errAddrNotFound
	}
	addressLinkedTokenContracts, err := s.meta.GetAddressTokens(addr.Index())

	result := resp.AddressTokensListResponse{}

	rpcProvider := s.getRPCProvider(addr.MHDA().Chain().Key())
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

	if !types.IsNetworkExists(dto.ChainKey) {
		return nil, errChainKeyNotSupported
	}

	defaultNodeIndex := s.getRPCProvider(dto.ChainKey).DefaultNodeId()

	ctx := types.NewRPCContext(dto.ChainKey, defaultNodeIndex)

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

	if !types.IsNetworkExists(dto.ChainKey) {
		return errChainKeyNotSupported
	}

	defaultNodeIndex := s.getRPCProvider(dto.ChainKey).DefaultNodeId()

	ctx := types.NewRPCContext(dto.ChainKey, defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return err
	}

	if provider.IsTokenConfigExists(dto.Contract) {
		tokenConfig = provider.GetTokenConfig(dto.Contract)

		tokenIndex = types.TokenIndex{
			ChainKey: dto.ChainKey,
			Contract: tokenConfig.Contract(),
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
			ChainKey: dto.ChainKey,
			Contract: tokenConfig.Contract(),
		}
	}

	s.meta.AddTokenConfig(dto.ChainKey, tokenConfig)

	if dto.MhdaPath != "" {
		addrKey, err := mhda.ParseNSS(dto.MhdaPath)
		if err != nil {
			return err
		}

		addr := s.meta.GetAddress(addrKey.NSS())

		if err != nil {
			return errors.New("address not found")
		}

		err = s.meta.SetTokenConfigAddressLink(addr.Index(), tokenIndex)
	}

	return err
}

func (s *Wallet) ExecuteRPC(dto *dto.ExecuteRPCRequestDTO) ([]byte, error) {
	if !types.IsNetworkExists(dto.ChainKey) {
		return nil, errChainKeyNotSupported
	}

	defaultNodeIndex := s.getRPCProvider(dto.ChainKey).DefaultNodeId()

	ctx := types.NewRPCContext(dto.ChainKey, defaultNodeIndex)

	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	return provider.GetBlock(ctx, 0)
}
