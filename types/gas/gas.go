package gas

type AlgorithmType string

type Calculator interface {
	SetPair(currency float64, suffix string)
	SuggestSlow() float64
	SuggestRegular() float64
	SuggestPriority() float64
	LimitMin() float64
	LimitMax() float64
	Format() string
}

type CalcOpts struct {
	GasSuffix     string
	TokenCurrency uint64
	TokenSuffix   string
	pairCurrency  float64
	pairSuffix    string
}

func NewGasCalculator(algorithm AlgorithmType, opts *CalcOpts) Calculator {
	switch algorithm {
	case AlgEVML1v1:
		return CalcEVML1V1{
			CalcOpts: opts,
		}
	case AlgBTCL1v1:
		return CalcBTCL1V1{
			CalcOpts: opts,
		}
	}
	return nil
}
