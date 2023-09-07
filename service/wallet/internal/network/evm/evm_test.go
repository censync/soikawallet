package evm

import "testing"

func TestEVM_floatToWei(t *testing.T) {
	var floatValues = []float64{122188888.8888, 9.999999, 0.33, 0.11, 0.22, 0.000005}
	for index := range floatValues {
		result := floatToWei(floatValues[index])
		t.Log(floatValues[index], result.String())
	}
}
