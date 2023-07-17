package gas

import (
	"encoding/json"
	"fmt"
)

const AlgEVML1v1 = AlgorithmType(`alg_evm_l1_1`)

type CalcEVML1V1 struct {
	*CalcOpts
	Units       float64 `json:"units"`
	BaseFee     float64 `json:"base_fee"`
	PriorityFee float64 `json:"priority_fee"`
	GasUsed     uint64  `json:"gas_used"`
	GasLimit    uint64  `json:"gas_limit"`
}

func NewCalcEVML1V1(calcOpts *CalcEVML1V1) Calculator {
	return calcOpts
}

func (c CalcEVML1V1) BaseGas() float64 {
	return c.BaseFee
}

func (c CalcEVML1V1) SuggestSlow() float64 {
	// return c.Units * (c.BaseFee) // low tip
	return c.PriorityFee // low tip
}

func (c CalcEVML1V1) SuggestRegular() float64 {
	//return c.Units * (c.BaseFee + c.PriorityFee*1.05) // suggest tip 5%
	return c.PriorityFee * 1.05 // suggest tip 5%
}

func (c CalcEVML1V1) SuggestPriority() float64 {
	/*if c.PriorityFee >= 1 {
		return c.Units * (c.BaseFee + c.PriorityFee*1.7) //  priority max 170%
	} else {
		return c.Units * (c.BaseFee + 1)
	}*/
	if c.PriorityFee >= float64(c.GasUnits) {
		return c.PriorityFee * 1.7 //  priority max 170%
	} else {
		return float64(c.GasUnits)
	}
}

func (c CalcEVML1V1) LimitMin() float64 {
	return 0 // ??
}

func (c CalcEVML1V1) LimitMax() uint64 {
	return c.GasLimit // max block ??
}

func (c CalcEVML1V1) LimitMaxGasFee(gasTipCap float64) float64 {
	return gasTipCap + c.BaseFee*2
}

func (c CalcEVML1V1) FormatHumanGas(gas float64) string {
	return fmt.Sprintf("%.3f", gas/float64(c.GasUnits))
}

func (c CalcEVML1V1) FormatHumanFiatPrice(gas float64) string {
	return fmt.Sprintf("%f$", gas/float64(c.GasUnits)*(c.FiatCurrency/float64(c.GasUnits)))
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
