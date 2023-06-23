package gas

const AlgEVML1v1 = AlgorithmType(`alg_evm_l1_1`)

type CalcEVML1V1 struct {
	*CalcOpts
	Units       float64
	BaseFee     float64
	PriorityFee float64
	MaxFee      float64
}

func (c CalcEVML1V1) SuggestSlow() float64 {
	//TODO implement me
	return c.Units * (c.BaseFee + 0.2) / 1e9
}

func (c CalcEVML1V1) SuggestRegular() float64 {
	//TODO implement me
	return c.Units * (c.BaseFee + c.PriorityFee) / 1e9
}

func (c CalcEVML1V1) SuggestPriority() float64 {
	//TODO implement me
	return c.Units * (c.BaseFee + c.PriorityFee) / 1e9
}

func (c CalcEVML1V1) LimitMin() float64 {
	//TODO implement me
	return c.Units * c.BaseFee / 1e9
}

func (c CalcEVML1V1) LimitMax() float64 {
	//TODO implement me
	return 30000 // max block
}

func (c CalcEVML1V1) Format() string {
	//TODO implement me
	return "{gas}"
}

func (c CalcEVML1V1) SetPair(currency float64, suffix string) {

}
