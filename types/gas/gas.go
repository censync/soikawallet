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

package gas

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AlgorithmType string

type Calculator interface {
	EstimateGas() uint64
	BaseGas() uint64
	SuggestSlow() uint64
	SuggestRegular() uint64
	SuggestPriority() uint64
	LimitMax() uint64
	Debug() string
	LimitMaxGasFee(gasTipCap uint64) uint64
	FormatHumanGas(gas uint64) string
	FormatHumanFiatPrice(gas uint64) string
	Marshal() ([]byte, error)
}

type CalcOpts struct {
	// GasEstimate is a count of gas, required for transaction
	GasEstimate uint64 `json:"gas_estimate"`
	// GasSymbol is a network defined gas unit name, e.g. "gwei", "satoshi"
	GasSymbol string `json:"gas_symbol"`
	// GasUnits is a units count per one base token, for evm = 1e9
	GasUnits uint64 `json:"gas_units"`
	// FiatSymbol is a short fiat currency suffix, e.g. $, €
	FiatSymbol string `json:"token_suffix"`
	// FiatCurrency is a fiat currency per one base token
	FiatCurrency float64 `json:"fiat_currency"`
}

/*
func (o CalcOpts) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}*/

/*
func NewGasCalculator(algorithm AlgorithmType, opts *CalcOpts) Calculator {
	switch algorithm {
	case AlgEVML1V1:
		return CalcEVML1V1{
			CalcOpts: opts,
		}
	case AlgBTCL1v1:
		return CalcBTCL1V1{
			CalcOpts: opts,
		}
	}
	return nil
}
*/

func Unmarshal(data []byte) (instance Calculator, err error) {
	var (
		tmp = struct {
			Type   AlgorithmType `json:"alg"`
			Config interface{}   `json:"config"`
		}{}
	)

	err = json.Unmarshal(data, &tmp)

	if err != nil {
		return nil, err
	}

	// TODO: Optimize serialization
	tmpCfg, _ := json.Marshal(tmp.Config)

	switch tmp.Type {
	case AlgEVML1V1:
		instance = &CalcEVML1V1{}
		err = json.Unmarshal(tmpCfg, instance)
	case AlgBTCL1v1:
		instance = &CalcBTCL1V1{}
		err = json.Unmarshal(tmpCfg, instance)
	default:
		err = errors.New(fmt.Sprintf(`undefined calculator algorithm %s`, tmp.Type))
	}

	return instance, err
}
