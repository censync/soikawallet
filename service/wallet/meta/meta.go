package meta

import (
	"encoding/json"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/types"
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
	// addresses key: mhda_nss
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

func (m *Meta) IsAddressExist(addrKey string) bool {
	_, ok := m.addresses[addrKey]
	return ok
}

func (m *Meta) Addresses() map[string]*Address {
	return m.addresses
}

func (m *Meta) GetAddress(addrKey string) *Address {
	return m.addresses[addrKey]
}

func (m *Meta) SetAddress(addrKey string, address *Address) {
	m.addresses[addrKey] = address
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
