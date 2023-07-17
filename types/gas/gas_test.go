package gas

import "testing"

const testData = `{"alg":"alg_evm_l1_1","config":{"gas_symbol":"gwei","gas_units":10000000000,"token_suffix":"$","fiat_currency":1900,"units":2100,"base_fee":24,"priority_fee":34,"gas_used":123123,"gas_limit":30000000}}`

func TestGas_Unmarshal(t *testing.T) {
	testCalcInstance, err := Unmarshal([]byte(testData))
	t.Log(testCalcInstance)
	t.Log(err)
}
