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
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"errors"
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/config/chain"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/censync/soikawallet/types/gas"
)

var (
	errTokenNotConfigured        = errors.New("token not configured")
	errTokenAllowanceApproveBase = errors.New("allowance base token")
	errTokenIncorrectType        = errors.New("incorrect token type")
	errUndefinedOperation        = errors.New("undefined operation")
)

func (s *Wallet) GetGasCalculatorConfig(dto *dto.GetGasCalculatorConfigDTO) (*resp.CalculatorConfig, error) {
	var (
		gasCalculator gas.Calculator
		gasConfig     map[string]uint64
		fiatSuffix    string
		fiatCurrency  float64
	)

	addrKey, err := mhda.ParseNSS(dto.MhdaPath)
	if err != nil {
		return nil, err
	}

	addr := s.meta.GetAddress(addrKey.NSS())

	if addr == nil {
		return nil, err
	}

	ctx := types.NewRPCContext(addr.MHDA().Chain().Key(), addr.NodeIndex(), addr.Address())
	provider, err := s.getNetworkProvider(ctx)

	if fiatPair := s.currenciesFiat.Get(provider.Currency()); fiatPair != nil {
		fiatSuffix = s.currenciesFiat.Symbol()
		fiatCurrency = fiatPair.Value()
	}

	if err != nil {
		return nil, err
	}

	// TODO: Optimize method
	if dto.Operation == "transfer" {
		if types.TokenStandard(dto.Standard) == types.TokenBase {
			gasConfig, err = provider.GetGasConfig(ctx)
		} else {
			tokenConfig := provider.GetTokenConfig(dto.Contract)

			if tokenConfig == nil {
				return nil, errTokenNotConfigured
			}
			gasConfig, err = provider.GetGasConfig(ctx, "transfer", dto.To, dto.Value, tokenConfig)
		}
	} else if dto.Operation == "approve" {
		if types.TokenStandard(dto.Standard) != types.TokenBase {
			tokenConfig := provider.GetTokenConfig(dto.Contract)

			if tokenConfig == nil {
				return nil, errTokenNotConfigured
			}

			gasConfig, err = provider.GetGasConfig(ctx, "approve", dto.To, dto.Value, tokenConfig)
		} else {
			return nil, errTokenAllowanceApproveBase
		}

	} else {
		return nil, errUndefinedOperation
	}

	if err != nil {
		return nil, err
	}

	switch addr.MHDA().Chain().Key() {
	case chain.EthereumChain.Key(), chain.PolygonChain.Key(), chain.Moonbeam.Key(), chain.MantleChain.Key():

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

	case chain.AvalancheCChain.Key():
		gasCalculator = gas.NewCalcEVML1V1(&gas.CalcEVML1V1{
			CalcOpts: &gas.CalcOpts{
				GasEstimate:  gasConfig["units"],
				GasSymbol:    "nAVAX",
				GasUnits:     1e9,
				FiatSymbol:   fiatSuffix,
				FiatCurrency: fiatCurrency,
			},
			BaseFee:     gasConfig["base_fee"],
			PriorityFee: gasConfig["priority_fee"],
			GasUsed:     gasConfig["gas_used"],
			GasLimit:    gasConfig["gas_limit"], // 30000 or 30e6?
		})
	// TODO: Add algorithms to init_wallet NetworkType
	case chain.BinanceSmartChain.Key(), chain.ArbitrumChain.Key(), chain.OptimismChain.Key(), chain.Linea.Key(), chain.BaseChain.Key():

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
		return nil, errors.New(fmt.Sprintf("gas calculator for network (%s) is not defined", addr.MHDA().Chain().Key()))
	}

	response := &resp.CalculatorConfig{}

	response.Calculator, err = gasCalculator.Marshal()

	if err != nil {
		return nil, err
	}

	return response, nil
}
