package gas

import "encoding/json"

const AlgBTCL1v1 = AlgorithmType(`alg_btc_l1_1`)

type CalcBTCL1V1 struct {
	*CalcOpts
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
	return 10e8
}

func (c CalcBTCL1V1) Format() string {
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
