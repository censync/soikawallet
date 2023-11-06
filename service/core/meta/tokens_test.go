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

package meta

import (
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	metaTokens        *tokens
	testTokenDecimals = 18
	testDataTokens    = [][]string{
		{
			"Test Token 1",
			"TEST1",
			"0x0000000000000000000000000000000000000000",
		},
		{
			"Test Token 2",
			"TEST2",
			"0x0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f",
		},
		{
			"Test Token 3",
			"TEST3",
			"0xffffffffffffffffffffffffffffffffffffffff",
		},
	}
	//chainKey = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1`)
)

func init() {
	metaTokens = &tokens{}
	metaTokens.initTokens()
}

func TestTokens_AddTokenConfig_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.subIndex)
	assert.NotNil(t, metaTokens.links)

	for index := range testDataTokens {
		tokenConfig := types.NewTokenConfig(
			types.TokenERC20,
			testDataTokens[index][0],
			testDataTokens[index][1],
			testDataTokens[index][2],
			testTokenDecimals,
		)

		err := metaTokens.AddTokenConfig(chainKey.Key(), tokenConfig)

		assert.Nil(t, err)
	}

	if len(metaTokens.tokens) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	if len(metaTokens.subIndex) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	for tokenIndex, tokenConfig := range metaTokens.tokens {
		assert.Equal(t, types.TokenERC20, tokenConfig.Standard())
		assert.Equal(t, testDataTokens[tokenIndex-1][0], tokenConfig.Name())
		assert.Equal(t, testDataTokens[tokenIndex-1][1], tokenConfig.Symbol())
		assert.Equal(t, testDataTokens[tokenIndex-1][2], tokenConfig.Contract())
		assert.Equal(t, testTokenDecimals, tokenConfig.Decimals())
	}

}

/*
func TestTokens_SetTokenConfigAddressLink_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.subIndex)
	assert.NotNil(t, metaTokens.links)

	for _, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			NetworkType: chainKey.key(),
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			Index:      1, // uint32(index) + 1,
			IsHardened: true,
		}
		err := metaTokens.SetTokenConfigAddressLink(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.Nil(t, err)
	}
}

func TestTokens_MarshalJSON_Positive(t *testing.T) {
	TestTokens_AddTokenConfig_Positive(t)
	TestTokens_SetTokenConfigAddressLink_Positive(t)
	str, err := json.Marshal(metaTokens)
	assert.Nil(t, err)
	fmt.Println(string(str))
}

func TestTokens_SetTokenConfigAddressLink_Negative_Duplicate(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	for index, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			CoinType: testNetwork,
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			InternalIndex: uint32(index) + 1,
			IsHardened:    true,
		}
		err := metaTokens.SetTokenConfigAddressLink(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.NotNil(t, err)
	}
}

func TestTokens_RemoveTokenConfigAddressLink_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	for index, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			CoinType: testNetwork,
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			InternalIndex: uint32(index) + 1,
			IsHardened:    true,
		}
		exists, err := metaTokens.IsTokenConfigAddressLinkExists(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.Nil(t, err)
		assert.Equal(t, true, exists)

		metaTokens.RemoveTokenConfigAddressLink(tokenIndex, types.AccountIndex(0), addressIndex)

		notExists, err := metaTokens.IsTokenConfigAddressLinkExists(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.Nil(t, err)
		assert.Equal(t, false, notExists)
	}
}

func TestTokens_RemoveTokenConfig_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	if len(metaTokens.tokens) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	if len(metaTokens.addressesLinks) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	for _, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			CoinType: testNetwork,
			Contract: entry[2],
		}

		err := metaTokens.RemoveTokenConfig(tokenIndex)
		assert.Nil(t, err)

		_, exists := metaTokens.tokens[tokenIndex]
		assert.False(t, exists)

		// _, IsLabelExists = metaTokens.addressesLinks[tokenIndex]
		// assert.False(t, IsLabelExists)
	}

	if len(metaTokens.tokens) > 0 {
		t.Fatal("awaiting zero length")
	}

	if len(metaTokens.addressesLinks) > 0 {
		t.Fatal("awaiting zero length")
	}
}
*/

func TestTokens_MarshalJSON(t *testing.T) {
	data, err := metaTokens.MarshalJSON()
	t.Log(data, err)
}
