package wallet

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	airgap "github.com/censync/go-airgap"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
	"strings"
)

func (s *Wallet) GetAllowance(dto *dto.GetTokenAllowanceDTO) (uint64, error) {
	dto.To = strings.TrimSpace(dto.To)
	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return 0, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return 0, err
	}

	if len(dto.To) < 4 {
		return 0, errors.New("incorrect recipient address")
	}

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return 0, err
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		return 0, errors.New("allowance not available for base tokens")
	}
	tokenConfig := provider.GetTokenConfig(dto.Contract)

	if tokenConfig == nil {
		return 0, errors.New("token not configured")
	}

	if tokenConfig.Standard() != types.TokenStandard(dto.Standard) {
		return 0, errors.New("incorrect token type")
	}
	return provider.GetTokenAllowance(ctx, tokenConfig.Contract(), dto.To)
}
func (s *Wallet) ApproveTokens(dto *dto.SendTokensDTO) (txId string, err error) {
	dto.To = strings.TrimSpace(dto.To)
	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return ``, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return ``, err
	}

	if addr.key == nil {
		return ``, errors.New("empty key for sign, use airgap option")
	}

	if len(dto.To) < 4 {
		return ``, errors.New("incorrect recipient address")
	}

	if dto.Contract == "" {
		return "", errors.New("contract not set")
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		return "", errors.New("cannot approve for base token")
	}

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return "", err
	}

	tokenConfig := provider.GetTokenConfig(dto.Contract)

	if tokenConfig == nil {
		return ``, errors.New("token not configured")
	}

	result, err := provider.TxApproveToken(ctx, dto.To, dto.Value, tokenConfig, dto.Gas, dto.GasTipCap, dto.GasFeeCap, addr.key.Get())

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
		Path: dto.DerivationPath,
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
	var addrKey *ecdsa.PrivateKey

	defer func() {
		addrKey = nil
	}()

	dto.To = strings.TrimSpace(dto.To)
	addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return nil, err
	}

	if !isAirGap {
		if addr.key == nil {
			return nil, errors.New("empty key for sign, use airgap option")
		}
		addrKey = addr.key.Get()
	}

	if len(dto.To) < 4 {
		return nil, errors.New("incorrect recipient address")
	}

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		isLegacy := false
		if addr.Network() != types.BSC {
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
			addrKey,
		)

	} else {
		tokenConfig := provider.GetTokenConfig(dto.Contract)

		if tokenConfig == nil {
			return nil, errors.New("token not configured")
		}

		if tokenConfig.Standard() != types.TokenStandard(dto.Standard) {
			return nil, errors.New("incorrect token type")
		}

		return provider.TxSendToken(
			ctx,
			dto.To,
			dto.Value,
			tokenConfig,
			dto.Gas,
			dto.GasTipCap,
			dto.GasFeeCap,
			addrKey,
		)
	}
}
