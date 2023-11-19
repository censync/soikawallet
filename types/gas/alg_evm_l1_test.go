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

import (
	"testing"
)

func TestCalcEVML1V1_MarshalJSON(t *testing.T) {
	gasInstance := NewCalcEVML1V1(&CalcEVML1V1{
		CalcOpts: &CalcOpts{
			GasSymbol:    "gwei",
			GasUnits:     1e9,
			FiatSymbol:   "$",
			FiatCurrency: 1900,
		},
		BaseFee:     24,
		PriorityFee: 34,
		GasUsed:     123123,
		GasLimit:    30000000,
	})

	data, err := gasInstance.Marshal()
	t.Log(string(data))
	t.Log(err)
}
