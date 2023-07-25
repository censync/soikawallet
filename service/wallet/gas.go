package wallet

import (
	"errors"
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/gas"
)

func (s *Wallet) GetGasCalculatorConfig(dto *dto.GetGasCalculatorConfigDTO) (*resp.CalculatorConfig, error) {
	var (
		gasCalculator gas.Calculator
		gasConfig     map[string]uint64
		fiatSuffix    string
		fiatCurrency  float64
	)

	addrPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addr, err := s.address(addrPath)

	if err != nil {
		return nil, err
	}

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if fiatPair := s.currenciesFiat.Get(provider.Currency()); fiatPair != nil {
		fiatSuffix = s.currenciesFiat.Symbol()
		fiatCurrency = fiatPair.Value()
	}

	// TODO: Optimize method
	if dto.Operation == "transfer" {
		if types.TokenStandard(dto.Standard) == types.TokenBase {
			gasConfig, err = provider.GetGasConfig(ctx)
		} else {
			tokenConfig := provider.GetTokenConfig(dto.Contract)

			if tokenConfig == nil {
				return nil, errors.New("token not configured")
			}
			gasConfig, err = provider.GetGasConfig(ctx, "transfer(address,uint256)", dto.To, dto.Value, tokenConfig)
		}
	} else if dto.Operation == "approve" {
		if types.TokenStandard(dto.Standard) != types.TokenBase {
			tokenConfig := provider.GetTokenConfig(dto.Contract)

			if tokenConfig == nil {
				return nil, errors.New("token not configured")
			}

			gasConfig, err = provider.GetGasConfig(ctx, "approve(address,uint256)", dto.To, dto.Value, tokenConfig)
		} else {
			return nil, errors.New("cannot approve for base token")
		}

	} else {
		return nil, errors.New("undefined operation")
	}

	if err != nil {
		return nil, err
	}

	switch addr.Network() {
	case types.Ethereum, types.Polygon:

		gasCalculator = gas.NewCalcEVML1V1(&gas.CalcEVML1V1{
			CalcOpts: &gas.CalcOpts{
				GasEstimate:  gasConfig["units"],
				GasSymbol:    "gwei",
				GasUnits:     1e9,
				FiatSymbol:   fiatSuffix,
				FiatCurrency: fiatCurrency,
			},
			BaseFee:     gasConfig["base_fee"],
			PriorityFee: gasConfig["priority_fee"],
			GasUsed:     gasConfig["gas_used"],
			GasLimit:    gasConfig["gas_limit"], // 30000 or 30e6?
		})
	case types.BSC:

		gasCalculator = gas.NewCalcEVML1V1(&gas.CalcEVML1V1{
			CalcOpts: &gas.CalcOpts{
				GasEstimate:  gasConfig["units"],
				GasSymbol:    "gwei",
				GasUnits:     1e9,
				FiatSymbol:   fiatSuffix,
				FiatCurrency: fiatCurrency,
			},
			BaseFee:     gasConfig["base_fee"],
			PriorityFee: gasConfig["priority_fee"],
			GasUsed:     gasConfig["gas_used"],
			GasLimit:    gasConfig["gas_limit"], // 30000 or 30e6?
		})
	default:
		return nil, errors.New(fmt.Sprintf("gas calculator for network (%d) is not defined", addr.Network()))
	}

	response := &resp.CalculatorConfig{}

	response.Calculator, err = gasCalculator.Marshal()

	if err != nil {
		return nil, err
	}

	return response, nil
}
