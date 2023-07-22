package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
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

func (s *Wallet) SendTokens(dto *dto.SendTokensDTO) (txId string, err error) {
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

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return "", err
	}

	if types.TokenStandard(dto.Standard) == types.TokenBase {
		if addr.Network() != types.BSC {
			return provider.TxSendBase(ctx, dto.To, dto.Value, dto.GasTipCap, dto.GasFeeCap, addr.key.Get())
		} else {
			// TODO: Optimize to AlgEVMLegacyV1
			return provider.TxSendBaseLegacy(ctx, dto.To, dto.Value, dto.GasTipCap, addr.key.Get())
		}

	} else {
		tokenConfig := provider.GetTokenConfig(dto.Contract)

		if tokenConfig == nil {
			return ``, errors.New("token not configured")
		}

		if tokenConfig.Standard() != types.TokenStandard(dto.Standard) {
			return ``, errors.New("incorrect token type")
		}
		return provider.TxSendToken(ctx, dto.To, dto.Value, tokenConfig, addr.key.Get())
	}
}
