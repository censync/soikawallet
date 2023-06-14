package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/types"
	"sync"
)

type tokens struct {
	mu             sync.RWMutex
	tokens         map[types.TokenIndex]*TokenConfigMeta
	addressesLinks map[uint32]map[types.AccountIndex][]types.AddressIndex
}
type TokenConfigMeta struct {
	*types.TokenConfig
	Index uint32
}

func (t *tokens) initTokens() {
	t.tokens = map[types.TokenIndex]*TokenConfigMeta{}
	t.addressesLinks = map[uint32]map[types.AccountIndex][]types.AddressIndex{}
}

func (t *tokens) AddTokenConfig(coinType types.CoinType, config *types.TokenConfig) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !types.IsCoinExists(coinType) {
		return errors.New("coin type is not set")
	}

	var lastIndex uint32

	for _, tokenConfig := range t.tokens {
		if tokenConfig.Index > lastIndex {
			lastIndex = tokenConfig.Index
		}
	}

	lastIndex++

	tokenIndex := types.TokenIndex{
		CoinType: coinType,
		Contract: config.Contract(),
	}
	t.tokens[tokenIndex] = &TokenConfigMeta{
		TokenConfig: config,
		Index:       lastIndex,
	}
	if _, ok := t.addressesLinks[lastIndex]; !ok {
		t.addressesLinks[lastIndex] = map[types.AccountIndex][]types.AddressIndex{}
	}
	return nil
}

func (t *tokens) RemoveTokenConfig(index types.TokenIndex) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.tokens[index]; !ok {
		return errors.New("token is not IsLabelExists")
	}

	delete(t.addressesLinks, t.tokens[index].Index)
	delete(t.tokens, index)

	return nil
}

// addresses links

func (t *tokens) IsTokenConfigAddressLinkExists(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) (bool, error) {
	metaTokenConfig, ok := t.tokens[tokenIndex]
	if !ok {
		return false, errors.New("token is not IsLabelExists")
	}

	if t.addressesLinks[metaTokenConfig.Index] != nil {
		for _, index := range t.addressesLinks[metaTokenConfig.Index][accountIndex] {
			if index == addressIndex {
				return true, nil
			}
		}
	}

	return false, nil
}

func (t *tokens) GetTokenConfigAddressLinks(tokenIndex types.TokenIndex, accountIndex types.AccountIndex) ([]types.AddressIndex, error) {
	metaTokenConfig, ok := t.tokens[tokenIndex]

	if !ok {
		return nil, errors.New("token is not IsLabelExists")
	}

	return t.addressesLinks[metaTokenConfig.Index][accountIndex], nil
}

func (t *tokens) SetTokenConfigAddressLink(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	isExist, err := t.IsTokenConfigAddressLinkExists(tokenIndex, accountIndex, addressIndex)
	if err != nil {
		return err
	}

	if isExist {
		return errors.New("address already linked")
	}

	metaTokenConfig := t.tokens[tokenIndex]
	t.addressesLinks[metaTokenConfig.Index][accountIndex] = append(t.addressesLinks[metaTokenConfig.Index][accountIndex], addressIndex)
	return nil
}

// GetAddressTokens TODO: Add composite index, linked to address,
// includes labels, node, tokens and other links
func (t *tokens) GetAddressTokens(coinType types.CoinType, accountIndex types.AccountIndex, addressIndex types.AddressIndex) ([]*types.TokenConfig, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []*types.TokenConfig

	if !types.IsCoinExists(coinType) {
		return nil, errors.New("coin type is not set")
	}

	allContractsPerCoin := map[uint32]types.TokenIndex{}

	for tokenIndex := range t.tokens {
		if tokenIndex.CoinType == coinType {
			allContractsPerCoin[t.tokens[tokenIndex].Index] = tokenIndex
		}
	}

	if len(allContractsPerCoin) == 0 {
		return result, nil
	}

	for index, tokenIndex := range allContractsPerCoin {
		for _, linkedAddressIndex := range t.addressesLinks[index][accountIndex] {
			if linkedAddressIndex == addressIndex {
				result = append(result, t.tokens[tokenIndex].TokenConfig)
				break
			}
		}
	}

	return result, nil
}

func (t *tokens) RemoveTokenConfigAddressLink(tokenIndex types.TokenIndex, accountIndex types.AccountIndex, addressIndex types.AddressIndex) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	metaTokenConfig, ok := t.tokens[tokenIndex]

	if !ok {
		return errors.New("token is not IsLabelExists")
	}

	for index := range t.addressesLinks[metaTokenConfig.Index][accountIndex] {
		if t.addressesLinks[metaTokenConfig.Index][accountIndex][index] == addressIndex {
			t.addressesLinks[metaTokenConfig.Index][accountIndex] = append(t.addressesLinks[metaTokenConfig.Index][accountIndex][:index], t.addressesLinks[metaTokenConfig.Index][accountIndex][index+1:]...)
		}
	}
	return nil
}

func (t *tokens) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tokensExport := map[string]*types.TokenConfig{}
	for tokenIndex, token := range t.tokens {
		tokensExport[fmt.Sprintf("%d:%d", tokenIndex.CoinType, token.Index)] = token.TokenConfig
	}

	result := map[string]interface{}{
		"tokens": tokensExport,
		"links":  t.addressesLinks,
	}

	return json.Marshal(result)
}
