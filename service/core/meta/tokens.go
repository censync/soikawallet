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

package meta

import (
	"encoding/json"
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
	"sync"
)

type tokens struct {
	mu sync.RWMutex

	tokens map[uint32]*types.TokenConfig

	// addressesLinks represents map[token_index]internal_map_enum_index
	subIndex map[types.TokenIndex]uint32
	links    map[aIndex][]uint32
}

func (t *tokens) initTokens() {
	t.tokens = map[uint32]*types.TokenConfig{}
	t.subIndex = map[types.TokenIndex]uint32{}
	t.links = map[aIndex][]uint32{}
	//t.addressesLinks = map[uint32]map[types.AccountIndex][]types.AddressIndex{}
}

func (t *tokens) AddTokenConfig(chainKey mhda.ChainKey, config *types.TokenConfig) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !types.IsNetworkExists(chainKey) {
		return errors.New("network type is not set")
	}

	var lastIndex uint32

	for _, lastIndex = range t.subIndex {
	}

	lastIndex++

	tokenIndex := types.TokenIndex{
		ChainKey: chainKey,
		Contract: config.Contract(),
	}

	t.subIndex[tokenIndex] = lastIndex
	t.tokens[lastIndex] = config

	return nil
}

func (t *tokens) RemoveTokenConfig(index types.TokenIndex) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if internalIdx, ok := t.subIndex[index]; ok {
		delete(t.subIndex, index)
		delete(t.tokens, internalIdx)
		// TODO: Also remove links
	} else {
		return errors.New("token is not exists")
	}

	return nil
}

// addresses links

func (t *tokens) IsTokenConfigAddressLinkExists(addrIdx aIndex, tokenIndex types.TokenIndex) (bool, error) {
	internalIndex, ok := t.subIndex[tokenIndex]
	if !ok {
		return false, errors.New("token is not exist")
	}

	if _, ok = t.links[addrIdx]; ok {
		for _, index := range t.links[addrIdx] {
			if index == internalIndex {
				return true, nil
			}
		}
	}
	// No error: token index correct, but aIndex isn't exists
	return false, nil
}

/*
func (t *tokens) GetTokenConfigAddressLinks(tokenIndex types.TokenIndex, accountIndex types.AccountIndex) ([]types.AddressIndex, error) {
	metaTokenConfig, ok := t.tokens[tokenIndex]

	if !ok {
		return nil, errors.New("token is not IsLabelExists")
	}

	return t.addressesLinks[metaTokenConfig.InternalIndex][accountIndex], nil
}*/

func (t *tokens) SetTokenConfigAddressLink(addrIdx aIndex, tokenIndex types.TokenIndex) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	isExist, err := t.IsTokenConfigAddressLinkExists(addrIdx, tokenIndex)
	if err != nil {
		return err
	}

	if isExist {
		return errors.New("address already linked")
	}

	internalIndex, ok := t.subIndex[tokenIndex]
	if !ok {
		return errors.New("token is not exist")
	}

	t.links[addrIdx] = append(t.links[addrIdx], internalIndex)
	return nil
}

// GetAddressTokens TODO: Add composite index, linked to address,
// includes labels, node, tokens and other links
func (t *tokens) GetAddressTokens(addrIdx aIndex) ([]*types.TokenConfig, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []*types.TokenConfig

	for _, internalIndex := range t.links[addrIdx] {
		if tokenConfig, ok := t.tokens[internalIndex]; ok {
			result = append(result, tokenConfig)
		}
	}

	return result, nil
}

func (t *tokens) RemoveTokenConfigAddressLink(addrIdx aIndex, tokenIndex types.TokenIndex) error {
	var linkExists bool

	t.mu.Lock()
	defer t.mu.Unlock()

	internalIndex, ok := t.subIndex[tokenIndex]

	if !ok {
		return errors.New("token is not exist")
	}

	for index := range t.links[addrIdx] {
		if t.links[addrIdx][index] == internalIndex {
			linkExists = true
			t.links[addrIdx] = append(t.links[addrIdx][:index], t.links[addrIdx][index+1:]...)
		}
	}

	if linkExists {
		return errors.New("link not exists")
	}

	return nil
}

func (t *tokens) MarshalJSON() ([]byte, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	/*tokensExport := map[uint32]*types.TokenConfig{}
	for subIdx, tokenConfig := range t.tokens {
		tokensExport[subIdx] = tokenConfig
	}*/

	result := map[string]interface{}{
		"tokens":    t.tokens,
		"sub_index": t.subIndex,
		"links":     t.links,
	}

	return json.Marshal(result)

}
