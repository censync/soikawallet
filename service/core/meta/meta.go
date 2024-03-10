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
	"encoding/json"
	"errors"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types"
	"sync/atomic"
	"time"
)

const (
	metaSettingsVersion = 2
)

// Meta structure contains labels for synchronization
// all user configuration with AirGap

type Meta struct {
	version uint8

	// nonce is the number of operations with meta config objects,
	// required for synchronization with AirGap Vault
	nonce          uint32
	nonceUpdatedAt int64 // UTC
	// addresses key: mhda_nss of address => val: address
	addresses map[string]*Address
	labels
	nodes
	tokens

	//deliveredKeys  []string
	//w3Accounts     []string
}

func InitMeta() *Meta {
	instance := &Meta{
		version: metaSettingsVersion,
		// debug
		//deliveredKeys:  []string{},
		addresses:      map[string]*Address{},
		nonce:          0,
		nonceUpdatedAt: time.Now().UTC().Unix(),
	}

	instance.initLabels()

	instance.initNodes()

	instance.initTokens()

	return instance
}

func (m *Meta) NonceAdd() {
	atomic.AddUint32(&m.nonce, 1)
	m.nonceUpdatedAt = time.Now().UTC().Unix()
}

func (m *Meta) IsAddressExist(nssKey string) bool {
	_, ok := m.addresses[nssKey]
	return ok
}

func (m *Meta) Addresses() map[string]*Address {
	return m.addresses
}

// GetAddress returns address mhda nss key
func (m *Meta) GetAddress(nssKey string) *Address {
	return m.addresses[nssKey]
}

func (m *Meta) SetAddress(nssKey string, address *Address) {
	m.addresses[nssKey] = address
}

func (m *Meta) RemoveAddress(nssKey string) error {
	if _, ok := m.addresses[nssKey]; !ok {
		return errors.New("address not exists")

	}
	m.addresses[nssKey].key.Free()
	delete(m.addresses, nssKey)
	return nil
}

// Nodes operations

func (m *Meta) AddRPCNode(index types.NodeIndex, rpc *types.RPC) error {
	err := m.nodes.AddRPCNode(index, rpc)
	if err == nil {
		m.NonceAdd()
	}
	return err
}
func (m *Meta) RemoveRPCNode(nodeIndex types.NodeIndex) error {
	err := m.nodes.RemoveRPCNode(nodeIndex)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) SetRPCAddressLink(addrIdx aIndex, nodeIndex types.NodeIndex) error {
	err := m.nodes.SetRPCAddressLink(addrIdx, nodeIndex)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) RemoveRPCAccountLink(addrIdx aIndex, nodeIndex types.NodeIndex) error {
	err := m.nodes.RemoveRPCAccountLink(addrIdx, nodeIndex)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

// Tokens

func (m *Meta) AddTokenConfig(chainKey mhda.ChainKey, config *types.TokenConfig) error {
	err := m.tokens.AddTokenConfig(chainKey, config)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) RemoveTokenConfig(index types.TokenIndex) error {
	err := m.tokens.RemoveTokenConfig(index)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) SetTokenConfigAddressLink(addrIdx aIndex, tokenIndex types.TokenIndex) error {
	err := m.tokens.SetTokenConfigAddressLink(addrIdx, tokenIndex)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) RemoveTokenConfigAddressLink(addrIdx aIndex, tokenIndex types.TokenIndex) error {
	err := m.tokens.RemoveTokenConfigAddressLink(addrIdx, tokenIndex)
	if err == nil {
		m.NonceAdd()
	}
	return err
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version      uint8  `json:"v"`
		Nonce        uint32 `json:"nonce"`
		NonceUpdated int64  `json:"nonce_ts"`
		Labels       labels `json:"labels"`
		Nodes        nodes  `json:"nodes"`
		Tokens       tokens `json:"tokens"`
	}{
		Version: m.version,
		Nonce:   m.nonce,
		Labels:  m.labels,
		Nodes:   m.nodes,
		Tokens:  m.tokens,
	})
}

func (m *Meta) UnmarshalJSON(b []byte) error {
	return nil
}
