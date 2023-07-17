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

func (c CalcEVML1V1) SuggestSlow() float64 {
	return c.Units * (c.BaseFee) // low tip
}

func (c CalcEVML1V1) SuggestRegular() float64 {
	return c.Units * (c.BaseFee + c.PriorityFee*1.05) // suggest tip 5%
}

func (c CalcEVML1V1) SuggestPriority() float64 {
	if c.PriorityFee >= 1 {
		return c.Units * (c.BaseFee + c.PriorityFee*1.7) //  priority max 170%
	} else {
		return c.Units * (c.BaseFee + 1)
	}

}

func (c CalcEVML1V1) LimitMin() float64 {
	return c.Units * c.BaseFee / float64(c.GasUnits) // ??
}

func (c CalcEVML1V1) LimitMax() uint64 {
	return c.GasLimit // max block ??
}

func (c CalcEVML1V1) Format() string {
	return fmt.Sprintf("Filled: %.1f%%, Slow: %f$, Regular: %f$, Priority: %f$,",
		float64(c.GasUsed)/float64(c.GasLimit)*100,
		c.SuggestSlow()/float64(c.GasUnits)*c.FiatCurrency,
		c.SuggestRegular()/float64(c.GasUnits)*c.FiatCurrency,
		c.SuggestPriority()/float64(c.GasUnits)*c.FiatCurrency,
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
