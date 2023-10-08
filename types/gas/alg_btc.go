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
)

const AlgBTCL1v1 = AlgorithmType(`alg_btc_l1_1`)

type CalcBTCL1V1 struct {
	*CalcOpts
}

func (c CalcBTCL1V1) EstimateGas() uint64 {
	return 0
}

func (c CalcBTCL1V1) BaseGas() uint64 {
	return 0
}

func (c CalcBTCL1V1) SuggestSlow() uint64 {
	//TODO implement me
	return 4770
}

func (c CalcBTCL1V1) SuggestRegular() uint64 {
	//TODO implement me
	return 5000
}

func (c CalcBTCL1V1) SuggestPriority() uint64 {
	//TODO implement me
	return 5830
}

func (c CalcBTCL1V1) LimitMax() uint64 {
	//TODO implement me
	return 1e8
}

func (c CalcBTCL1V1) LimitMaxGasFee(gasTipCap uint64) uint64 {
	return 1
}

func (c CalcBTCL1V1) FormatHumanGas(gas uint64) string {
	return ""
}

func (c CalcBTCL1V1) FormatHumanFiatPrice(gas uint64) string {
	return ""
}

func (c CalcBTCL1V1) Debug() string {
	//TODO implement me
	return "{gas}"
}

func (c CalcBTCL1V1) Marshal() ([]byte, error) {
	var export = struct {
		Type   AlgorithmType `json:"alg"`
		Config *CalcBTCL1V1  `json:"config"`
	}{
		Type:   AlgBTCL1v1,
		Config: &c,
	}
	return json.Marshal(export)
}
