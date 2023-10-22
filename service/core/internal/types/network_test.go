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

package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	testRPCTitle     = "My test RPC"
	testRPCEndpoint  = "https://rpc.example.com/testnet"
	testRPCIsDefault = true
)

func TestRPC_MarshalJSON(t *testing.T) {
	rpc := NewRPC(
		testRPCTitle,
		testRPCEndpoint,
		testRPCIsDefault,
	)

	strJSON, err := json.Marshal(rpc)
	assert.Nil(t, err)

	assert.Equal(t, []byte(fmt.Sprintf(
		`["%s","%s","%s"]`,
		testRPCTitle,
		testRPCEndpoint,
		strconv.FormatBool(testRPCIsDefault)),
	),
		strJSON,
	)
}
