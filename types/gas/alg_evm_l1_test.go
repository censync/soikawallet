package gas

import (
	"testing"
)

func TestCalcEVML1V1_MarshalJSON(t *testing.T) {
	gasInstance := NewCalcEVML1V1(&CalcEVML1V1{
		CalcOpts: &CalcOpts{
			GasSymbol:    "gwei",
			GasUnits:     10e9,
			TokenSuffix:  "$",
			FiatCurrency: 1900,
		},
		Units:       2100,
		BaseFee:     24,
		PriorityFee: 34,
		GasUsed:     123123,
		GasLimit:    30000000,
	})

	data, err := gasInstance.Marshal()
	t.Log(string(data))
	t.Log(err)
}
