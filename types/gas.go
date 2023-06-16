package types

const (
	GasAlgEVML1v1 = AlgorithmType(`alg_evm_l1_1`)
	GasAlgEVML2   = AlgorithmType(`alg_evm_l2_1`)
	GasAlgTronV1  = AlgorithmType(`alg_evm_tron_1`)
)

type AlgorithmType string

type GasCalcOpts struct {
	GasCalculator
	GasSuffix     string  // gwei
	TokenCurrency uint64  // 10e8
	TokenSuffix   string  // Eth
	FiatSuffix    string  // USD
	FiatCurrency  float64 // 1700
}

type GasCalculator interface {
	SuggestSlow()
	SuggestRegular()
	SuggestPriority()
	LimitMin()
	LimitMax()
	FormatGas()
	FormatToken()
	FormatCurrency()
}

type CalcBTC struct {
	*GasCalcOpts
}

type CalcEVML1V1 struct {
	*GasCalcOpts
	Units       float64
	BaseFee     float64
	PriorityFee float64
}

type CalcEVML2 struct {
	*GasCalcOpts
}

type CalcTron struct {
	*GasCalcOpts
}

func NewGasCalculator(algorithm AlgorithmType, opts *GasCalcOpts) GasCalculator {
	switch algorithm {
	case GasAlgEVML1v1:
		return CalcEVML1V1{
			GasCalcOpts: opts,
		}
	}
	return nil
}
