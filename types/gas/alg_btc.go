package gas

import (
	"encoding/json"
	"fmt"
)

const AlgBTCL1v1 = AlgorithmType(`alg_btc_l1_1`)

type CalcBTCL1V1 struct {
	*CalcOpts
}

func (c CalcBTCL1V1) BaseGas() float64 {
	return 0
}

func (c CalcBTCL1V1) SuggestSlow() float64 {
	//TODO implement me
	return 4770
}

func (c CalcBTCL1V1) SuggestRegular() float64 {
	//TODO implement me
	return 5000
}

func (c CalcBTCL1V1) SuggestPriority() float64 {
	//TODO implement me
	return 5830
}

func (c CalcBTCL1V1) LimitMin() float64 {
	//TODO implement me
	return 4770
}

func (c CalcBTCL1V1) LimitMax() uint64 {
	//TODO implement me
	return 1e8
}

func (c CalcBTCL1V1) LimitMaxGasFee(gasTipCap float64) float64 {
	return 1
}

func (c CalcBTCL1V1) FormatHumanGas(gas float64) string {
	return fmt.Sprintf("%.3f", gas/float64(c.GasUnits))
}

func (c CalcBTCL1V1) FormatHumanFiatPrice(gas float64) string {
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
