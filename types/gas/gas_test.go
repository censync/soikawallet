// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package gas

import "testing"

const testData = `{"alg":"alg_evm_l1_1","config":{"gas_symbol":"gwei","gas_units":10000000000,"token_suffix":"$","fiat_currency":1900,"units":2100,"base_fee":24,"priority_fee":34,"gas_used":123123,"gas_limit":30000000}}`

func TestGas_Unmarshal(t *testing.T) {
	testCalcInstance, err := Unmarshal([]byte(testData))
	t.Log(testCalcInstance)
	t.Log(err)
}
