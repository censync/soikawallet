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
	"fmt"
)

const AlgEVML1V1 = AlgorithmType(`alg_evm_l1_1`)

type CalcEVML1V1 struct {
	*CalcOpts
	BaseFee     uint64 `json:"base_fee"`
	PriorityFee uint64 `json:"priority_fee"`
	GasUsed     uint64 `json:"gas_used"`
	GasLimit    uint64 `json:"gas_limit"`
}

func NewCalcEVML1V1(calcOpts *CalcEVML1V1) Calculator {
	return calcOpts
}

func (c CalcEVML1V1) EstimateGas() uint64 {
	return c.GasEstimate
}

func (c CalcEVML1V1) BaseGas() uint64 {
	return c.BaseFee
}

func (c CalcEVML1V1) SuggestSlow() uint64 {
	return uint64(float64(c.PriorityFee) * 1.05) // low tip
}

func (c CalcEVML1V1) SuggestRegular() uint64 {
	return uint64(float64(c.PriorityFee) * 1.45) // suggest tip 45%
}

func (c CalcEVML1V1) SuggestPriority() uint64 {
	return uint64(float64(c.PriorityFee) * 1.7) // max tip 170%
}

func (c CalcEVML1V1) LimitMax() uint64 {
	return c.GasLimit // max block ??
}

func (c CalcEVML1V1) LimitMaxGasFee(gasTipCap uint64) uint64 {
	return gasTipCap + c.BaseFee*2
}

func (c CalcEVML1V1) FormatHumanGas(gas uint64) string {
	return fmt.Sprintf("%.3f", float64(gas)/float64(c.GasUnits))
}

func (c CalcEVML1V1) FormatHumanFiatPrice(gas uint64) string {
	return fmt.Sprintf("%f$", float64(gas)/float64(c.GasUnits)*(c.FiatCurrency/float64(c.GasUnits)))
}

func (c CalcEVML1V1) Debug() string {
	return fmt.Sprintf("Filled block: %.1f%%",
		float64(c.GasUsed)/float64(c.GasLimit)*100,
	)
}

// Marshal instead MarshalJSON, for deprecation recursively calling inner MarshalJSON in shadow struct
func (c CalcEVML1V1) Marshal() ([]byte, error) {
	var export = struct {
		Type   AlgorithmType `json:"alg"`
		Config *CalcEVML1V1  `json:"config"`
	}{
		Type:   AlgEVML1V1,
		Config: &c,
	}
	return json.Marshal(&export)
}
