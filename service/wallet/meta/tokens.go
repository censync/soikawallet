package meta

import (
	"errors"
	"github.com/censync/soikawallet/types"
	"sync"
)

type tokens struct {
	mu             sync.RWMutex
	tokens         map[types.TokenIndex]*types.TokenConfig
	addressesLinks map[types.TokenIndex]map[types.AccountIndex][]types.AddressIndex
}

func (t *tokens) initTokens() {
	t.tokens = map[types.TokenIndex]*types.TokenConfig{}
	t.addressesLinks = map[types.TokenIndex]map[types.AccountIndex][]types.AddressIndex{}
}

func (t *tokens) AddTokenConfig(index types.TokenIndex, config *types.TokenConfig) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	t.tokens[index] = config
	if _, ok := t.addressesLinks[index]; !ok {
		t.addressesLinks[index] = map[types.AccountIndex][]types.AddressIndex{}
	}
}

func (t *tokens) RemoveTokenConfig(index types.TokenIndex) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.tokens[index]; !ok {
		return errors.New("node is not exists")
	}

	delete(t.tokens, index)
	delete(t.addressesLinks, index)

	return nil
}

// addresses links

func (t *tokens) IsTokenConfigAddressLinkExists(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) bool {
	if t.addressesLinks[tokenIndex] != nil {
		for _, index := range t.addressesLinks[tokenIndex][accountIndex] {
			if index == addressIndex {
				return true
			}
		}
	}
	return false
}

func (t *tokens) GetTokenConfigAddressLinks(tokenIndex types.TokenIndex, accountIndex types.AccountIndex) []types.AddressIndex {
	return t.addressesLinks[tokenIndex][accountIndex]
}

func (t *tokens) SetTokenConfigAddressLink(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.IsTokenConfigAddressLinkExists(tokenIndex, accountIndex, addressIndex) {
		return errors.New("address already linked")
	}

	t.addressesLinks[tokenIndex][accountIndex] = append(t.addressesLinks[tokenIndex][accountIndex], addressIndex)
	return nil
}

func (t *tokens) GetTokenConfigAddressLinksCount(coinType types.CoinType, contract string) int {
	return len(t.addressesLinks[types.TokenIndex{
		CoinType: coinType,
		Contract: contract,
	}])
}

func (t *tokens) GetAddressTokens(coinType types.CoinType, accountIndex types.AccountIndex, addressIndex types.AddressIndex) ([]string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []string

	if !types.IsCoinExists(coinType) {
		return nil, errors.New("coin type is not set")
	}

	for tokenIndex := range t.addressesLinks {
		if _, ok := t.addressesLinks[tokenIndex][accountIndex]; ok {
			for _, addr := range t.addressesLinks[tokenIndex][accountIndex] {
				if addr == addressIndex {
					result = append(result, tokenIndex.Contract)
					break
				}
			}
		}
	}
	return result, nil
}

func (t *tokens) RemoveTokenConfigAddressLink(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for index := range t.addressesLinks[tokenIndex][accountIndex] {
		if t.addressesLinks[tokenIndex][accountIndex][index] == addressIndex {
			t.addressesLinks[tokenIndex][accountIndex] = append(t.addressesLinks[tokenIndex][accountIndex][:index], t.addressesLinks[tokenIndex][accountIndex][index+1:]...)
		}
	}
}
