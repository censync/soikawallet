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
	"github.com/censync/soikawallet/types"
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

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		isLegacy := false
		switch addr.MHDA().Chain().Key() {
		case chain.BinanceSmartChain.Key(), chain.ArbitrumChain.Key(), chain.OptimismChain.Key():
			isLegacy = true
		}

		return provider.TxSendBase(
			ctx,
			dto.To,
			dto.Value,
			dto.Gas,
			dto.GasTipCap,
			dto.GasFeeCap,
			isLegacy,
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

		return provider.TxSendToken(
			ctx,
			dto.To,
			dto.Value,
			tokenConfig,
			dto.Gas,
			dto.GasTipCap,
			dto.GasFeeCap,
			key,
		)
	}
}
