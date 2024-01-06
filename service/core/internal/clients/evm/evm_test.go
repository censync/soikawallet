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

var (
	testAmountStrValues = []string{
		"999", "999.10", "12345.54321",
		"122188888.888811112233344455", "9.12345", "0.123", "0.0000012345",
	}
	testAmountWeiValues = []string{
		"999000000000000000000", "999100000000000000000", "12345543210000000000000",
		"122188888888811112233344455", "9123450000000000000", "123000000000000000",
		"1234500000000",
	}
)

func TestEVM_strToWei(t *testing.T) {
	for index := range testAmountStrValues {
		result, err := strToWei(testAmountStrValues[index])

		if err != nil {
			t.Fatal(err)
		}

		if result.String() != testAmountWeiValues[index] {
			t.Fatalf("cannot convert value %x", testAmountStrValues[index])
		}
	}
}
