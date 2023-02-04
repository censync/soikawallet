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

func (n *tokens) initTokens() {
	n.tokens = map[types.TokenIndex]*types.TokenConfig{}
	n.addressesLinks = map[types.TokenIndex]map[types.AccountIndex][]types.AddressIndex{}
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
		if t.addressesLinks[tokenIndex][accountIndex] != nil {
			for _, index := range t.addressesLinks[tokenIndex][accountIndex] {
				if index == addressIndex {
					return true
				}
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
		return errors.New("already enabled")
	}

	t.addressesLinks[tokenIndex][accountIndex] = append(t.addressesLinks[tokenIndex][accountIndex], addressIndex)
	return nil
}

func (t *tokens) GetTokenConfigAddressLinksCount(coinType types.CoinType, tokenIndex uint32) int {
	return len(t.addressesLinks[types.TokenIndex{
		CoinType: coinType,
		Index:    tokenIndex,
	}])
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
