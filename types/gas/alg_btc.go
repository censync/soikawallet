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
