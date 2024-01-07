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
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	airgap "github.com/censync/go-airgap"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/config/chain"
	"github.com/censync/soikawallet/service/core/internal/clients/evm"
	"github.com/censync/soikawallet/service/core/internal/types"
	"strings"
)

var (
	errOpAddrRecipientIncorrect = errors.New("incorrect recipient address")
)

func (s *Wallet) GetAllowance(dto *dto.GetTokenAllowanceDTO) (uint64, error) {
	dto.To = strings.TrimSpace(dto.To)

	addrKey, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return 0, err
	}

	addr := s.meta.GetAddress(addrKey.NSS())

	if addr == nil {
		return 0, err
	}

	if len(dto.To) < 4 {
		return 0, errOpAddrRecipientIncorrect
	}

	ctx := types.NewRPCContext(addr.MHDA().Chain().Key(), addr.NodeIndex(), addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return 0, err
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		return 0, errTokenAllowanceApproveBase
	}
	tokenConfig := provider.GetTokenConfig(dto.Contract)

	if tokenConfig == nil {
		return 0, errTokenNotConfigured
	}

	if tokenConfig.Standard() != types.TokenStandard(dto.Standard) {
		return 0, errTokenNotConfigured
	}
	return provider.GetTokenAllowance(ctx, tokenConfig.Contract(), dto.To)
}
func (s *Wallet) ApproveTokens(dto *dto.SendTokensDTO) (txId string, err error) {
	dto.To = strings.TrimSpace(dto.To)
	addrKey, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return ``, err
	}

	addr := s.meta.GetAddress(addrKey.NSS())

	if addr == nil {
		return ``, err
	}

	if err != nil {
		return ``, err
	}

	if addr.Key() == nil {
		return ``, errors.New("empty key for sign, use airgap option")
	}

	if len(dto.To) < 4 {
		return ``, errOpAddrRecipientIncorrect
	}

	if dto.Contract == "" {
		return "", errors.New("contract not set")
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		return "", errTokenAllowanceApproveBase
	}

	ctx := types.NewRPCContext(addr.MHDA().Chain().Key(), addr.NodeIndex(), addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return "", err
	}

	tokenConfig := provider.GetTokenConfig(dto.Contract)

	if tokenConfig == nil {
		return ``, errTokenNotConfigured
	}

	result, err := provider.TxApproveToken(ctx, dto.To, dto.Value, tokenConfig, dto.Gas, dto.GasTipCap, dto.GasFeeCap, addr.Key().Get())

	if err != nil {
		return "", err
	}
	return result.(string), nil
}

type TxSendPrepare struct {
	Path string `json:"path"`
	Data []byte `json:"data"`
}

func (s *Wallet) SendTokensPrepare(dto *dto.SendTokensDTO) (*resp.AirGapMessage, error) {
	result, err := s.sendTokensProcess(dto, true)

	if err != nil {
		return nil, err
	}

	airGapMessage := &TxSendPrepare{
		Path: dto.MhdaPath,
		Data: result.([]byte),
	}

	airGapMessageBin, err := json.Marshal(&airGapMessage)

	if err != nil {
		return nil, errors.New("cannot marshal AirGap message")
	}

	airgapMsg := airgap.NewAirGap(airgap.VersionDefault, s.instanceId).
		CreateMessage().
		AddOperation(types.OpTxSend, airGapMessageBin)

	chunks, err := airgapMsg.MarshalB64Chunks()

	if err != nil {
		return nil, err
	}

	return &resp.AirGapMessage{
		Chunks: chunks,
	}, nil
}

func (s *Wallet) SendTokens(dto *dto.SendTokensDTO) (string, error) {
	result, err := s.sendTokensProcess(dto, false)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (s *Wallet) sendTokensProcess(dto *dto.SendTokensDTO, isAirGap bool) (interface{}, error) {
	var key *ecdsa.PrivateKey

	defer func() {
		key = nil
	}()

	dto.To = strings.TrimSpace(dto.To)
	addrKey, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return nil, err
	}

	addr := s.meta.GetAddress(addrKey.NSS())

	if addr == nil {
		return nil, err
	}

	if !isAirGap {
		if addr.Key() == nil {
			return nil, errors.New("empty key for sign, use airgap option")
		}
		key = addr.Key().Get()
	}

	if len(dto.To) < 4 {
		return nil, errOpAddrRecipientIncorrect
	}

	ctx := types.NewRPCContext(addr.MHDA().Chain().Key(), addr.NodeIndex(), addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	txFlag, err := txTypeByChainKey(addr.MHDA().Chain().Key())

	if err != nil {
		return nil, err
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		return provider.TxSendBase(
			ctx,
			dto.To,
			dto.Value,
			dto.Gas,
			dto.GasTipCap,
			dto.GasFeeCap,
			txFlag,
			key,
		)

	} else {
		tokenConfig := provider.GetTokenConfig(dto.Contract)

		if tokenConfig == nil {
			return nil, errTokenNotConfigured
		}

		if tokenConfig.Standard() != types.TokenStandard(dto.Standard) {
			return nil, errTokenIncorrectType
		}
		// TODO: Add txFlag
		return provider.TxSendToken(
			ctx,
			dto.To,
			dto.Value,
			tokenConfig,
			dto.Gas,
			dto.GasTipCap,
			dto.GasFeeCap,
			txFlag,
			key,
		)
	}
}

func txTypeByChainKey(chainKey mhda.ChainKey) (uint8, error) {
	txFlag := uint8(0)
	switch chainKey {
	case chain.EthereumChain.Key(), chain.PolygonChain.Key(), chain.Moonbeam.Key(), chain.AvalancheCChain.Key():
		txFlag = evm.TxFlagDynamic
	case chain.BinanceSmartChain.Key():
		txFlag = evm.TxFlagLegacy
	case chain.OptimismChain.Key(), chain.MantleChain.Key(), chain.ZkSyncEra.Key(): // chain.BaseChain.Key() test
		txFlag = evm.TxFlagL2
	// case chain.ArbitrumChain.Key(): implement
	default:
		return 0, errors.New("tx type not set for chain")
	}
	return txFlag, nil
}
