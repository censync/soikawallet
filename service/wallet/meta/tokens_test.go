package meta

import (
	"github.com/censync/soikawallet/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	metaTokens        *tokens
	testCoin          = types.Ethereum
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
)

func init() {
	metaTokens = &tokens{}
	metaTokens.initTokens()
}

func TestTokens_AddTokenConfig_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	for index := range testDataTokens {
		tokenConfig := types.NewTokenConfig(
			types.TokenERC20,
			testDataTokens[index][0],
			testDataTokens[index][1],
			testDataTokens[index][2],
			testTokenDecimals,
		)

		err := metaTokens.AddTokenConfig(testCoin, tokenConfig)

		assert.Nil(t, err)
	}

	if len(metaTokens.tokens) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	if len(metaTokens.addressesLinks) != len(testDataTokens) {
		t.Fatal("incorrect length")
	}

	/*
		for tokenIndex, tokenConfig := range metaTokens.tokens {
			index := tokenIndex.Index - 1
			assert.Equal(t, types.TokenERC20, tokenConfig.Standard())
			assert.Equal(t, testDataTokens[index][0], tokenConfig.Name())
			assert.Equal(t, testDataTokens[index][1], tokenConfig.Symbol())
			assert.Equal(t, testDataTokens[index][2], tokenConfig.Contract())
			assert.Equal(t, testTokenDecimals, tokenConfig.Decimals())
		}
	*/
}

func TestTokens_SetTokenConfigAddressLink_Positive(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	for index, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			CoinType: testCoin,
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			Index:      uint32(index) + 1,
			IsHardened: true,
		}
		err := metaTokens.SetTokenConfigAddressLink(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.Nil(t, err)
	}
}

func TestTokens_SetTokenConfigAddressLink_Negative_Duplicate(t *testing.T) {
	assert.NotNil(t, metaTokens)
	assert.NotNil(t, metaTokens.tokens)
	assert.NotNil(t, metaTokens.addressesLinks)

	for index, entry := range testDataTokens {
		tokenIndex := types.TokenIndex{
			CoinType: testCoin,
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			Index:      uint32(index) + 1,
			IsHardened: true,
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
			CoinType: testCoin,
			Contract: entry[2],
		}
		addressIndex := types.AddressIndex{
			Index:      uint32(index) + 1,
			IsHardened: true,
		}
		exists, err := metaTokens.IsTokenConfigAddressLinkExists(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.NotNil(t, err)
		assert.Equal(t, true, exists)

		metaTokens.RemoveTokenConfigAddressLink(tokenIndex, types.AccountIndex(0), addressIndex)

		notExists, err := metaTokens.IsTokenConfigAddressLinkExists(tokenIndex, types.AccountIndex(0), addressIndex)

		assert.NotNil(t, err)
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
			CoinType: testCoin,
			Contract: entry[2],
		}

		err := metaTokens.RemoveTokenConfig(tokenIndex)
		assert.Nil(t, err)

		_, exists := metaTokens.tokens[tokenIndex]
		assert.False(t, exists)

		// _, exists = metaTokens.addressesLinks[tokenIndex]
		// assert.False(t, exists)
	}

	if len(metaTokens.tokens) > 0 {
		t.Fatal("awaiting zero length")
	}

	if len(metaTokens.addressesLinks) > 0 {
		t.Fatal("awaiting zero length")
	}
}
