package meta

import (
	"encoding/json"
	"time"
)

const (
	metaSettingsVersion = 2
)

// Meta structure contains labels for synchronization
// all user configuration with AirGap

type Meta struct {
	version        uint8
	nonce          uint32
	nonceUpdatedAt uint64 // UTC
	// addresses key: mhda_nss
	addresses map[string]*Address
	//deliveredKeys  []string
	//w3Accounts     []string
	labels
	nodes
	tokens
}

func InitMeta() *Meta {
	instance := &Meta{
		version: metaSettingsVersion,
		// debug
		//deliveredKeys:  []string{},
		addresses:      map[string]*Address{},
		nonce:          0,
		nonceUpdatedAt: uint64(time.Now().UTC().Unix()),
	}

	instance.initLabels()

	instance.initNodes()

	instance.initTokens()

	return instance
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

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version       uint8    `json:"v"`
		Nonce         uint32   `json:"nonce"` // TODO: Add updated at
		DeliveredKeys []string `json:"delivered_keys"`
		//W3Accounts    []string `json:"w3_accounts"`
		Labels labels `json:"labels"`
		Nodes  nodes  `json:"nodes"`
		Tokens tokens `json:"tokens"`
	}{
		Version: m.version,
		Nonce:   m.nonce,
		//DeliveredKeys: m.deliveredKeys,
		//W3Accounts:    m.w3Accounts,
		Labels: m.labels,
		Nodes:  m.nodes,
		Tokens: m.tokens,
	})
}

func (m *Meta) UnmarshalJSON(b []byte) error {
	return nil
}
