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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"testing"
)

func TestEVM_strToWei(t *testing.T) {
	var floatValues = []string{"999", "999.0", "12345.54321", "122188888.888811112233344455", "9.999999", "0.33", "0.11", "0.22", "0.000005"}
	for index := range floatValues {
		result, err := strToWei(floatValues[index])

		if err != nil {
			t.Fatal(err)
		}

		t.Log(floatValues[index], result.String())

	}
}
