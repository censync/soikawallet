package gas

import (
	"encoding/json"
	"fmt"
)

const AlgEVML1v1 = AlgorithmType(`alg_evm_l1_1`)

type CalcEVML1V1 struct {
	*CalcOpts
	Units       uint64 `json:"units"`
	BaseFee     uint64 `json:"base_fee"`
	PriorityFee uint64 `json:"priority_fee"`
	GasUsed     uint64 `json:"gas_used"`
	GasLimit    uint64 `json:"gas_limit"`
}

func NewCalcEVML1V1(calcOpts *CalcEVML1V1) Calculator {
	return calcOpts
}

func (c CalcEVML1V1) BaseGas() uint64 {
	return c.BaseFee
}

func (c CalcEVML1V1) SuggestSlow() uint64 {
	// return c.Units * (c.BaseFee) // low tip
	return uint64(float64(c.PriorityFee) * 1.05) // low tip
}

func (c CalcEVML1V1) SuggestRegular() uint64 {
	//return c.Units * (c.BaseFee + c.PriorityFee*1.05) // suggest tip 5%
	return uint64(float64(c.PriorityFee) * 1.55) // suggest tip 55%
}

func (c CalcEVML1V1) SuggestPriority() uint64 {
	/*if c.PriorityFee >= 1 {
		return c.Units * (c.BaseFee + c.PriorityFee*1.7) //  priority max 170%
	} else {
		return c.Units * (c.BaseFee + 1)
	}*/
	return uint64(float64(c.PriorityFee) * 1.8)
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
		Type:   AlgEVML1v1,
		Config: &c,
	}
	return json.Marshal(&export)
}
