package wallet

import (
	"errors"
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/gas"
)

// GetGasPriceBaseTx Deprecated
func (s *Wallet) GetGasPriceBaseTx(dto *dto.GetGasPriceBaseTxDTO) (map[string]float64, error) {
	/*addressPath, err := types.ParsePath(dto.DerivationPath)
	if err != nil {
		return nil, err
	}

	addr, err := s.address(addressPath)

	if err != nil {
		return nil, err
	}

	ctx := types.NewRPCContext(addr.Network(), addr.nodeIndex, addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if err != nil {
		return nil, err
	}
	*/
	return nil, nil // provider.GetGasBaseTx(ctx)
}

func (s *Wallet) GetGasCalculatorConfig(dto *dto.GetAddressCalculatorConfigDTO) (*resp.CalculatorConfig, error) {
	var (
		gasCalculator gas.Calculator
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

	switch addr.Network() {
	case types.Ethereum:
		gasConfig, err := provider.GetGasBaseTx(ctx)

		if err != nil {
			return nil, err
		}

		gasCalculator = gas.NewCalcEVML1V1(&gas.CalcEVML1V1{
			CalcOpts: &gas.CalcOpts{
				GasSymbol:    "gwei",
				GasUnits:     1e9,
				FiatSymbol:   fiatSuffix,
				FiatCurrency: fiatCurrency,
			},
			Units:       gasConfig["units"],
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
