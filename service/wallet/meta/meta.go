package meta

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	metaSettingsVersion = 1
)

// Meta structure contains data for synchronization
// all user configuration with AirGap

type Meta struct {
	version        uint8
	nonce          uint32
	nonceUpdatedAt uint64 // UTC
	deliveredKeys  []string
	w3Accounts     []string
	labels
	nodes
	tokens
}

func InitMeta() *Meta {
	instance := &Meta{
		version: metaSettingsVersion,
		// debug
		deliveredKeys:  []string{"m/44'/60'/130'/0", "m/44'/60'/130'/1", "m/44'/60'/130'/2"},
		nonce:          42,
		nonceUpdatedAt: uint64(time.Now().UTC().Unix()),
	}

	instance.initLabels()

	instance.initNodes()

	instance.initTokens()

	return instance
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version       uint8    `json:"v"`
		Nonce         uint32   `json:"nonce"` // TODO: Add updated at
		DeliveredKeys []string `json:"delivered_keys"`
		Labels        labels   `json:"labels"`
		Nodes         nodes    `json:"nodes"`
		Tokens        tokens   `json:"tokens"`
	}{
		Version:       m.version,
		Nonce:         m.nonce,
		DeliveredKeys: m.deliveredKeys,
		Labels:        m.labels,
		Nodes:         m.nodes,
		Tokens:        m.tokens,
	})
}

func (m *Meta) UnmarshalJSON(b []byte) error {
	var result struct {
		Version uint8                           `json:"v"`
		Labels  map[LabelType]map[uint32]string `json:"labels"`
	}

	if err := json.Unmarshal(b, &result); err != nil {
		return err
	}

	*m = *InitMeta()

	if result.Version != m.version {
		if result.Version < m.version {
			return errors.New("config version is older version")
		} else {
			return errors.New("config version is newer version")
		}
	}

	for index, label := range result.Labels[AccountLabel] {
		m.labelsAccount.data[index] = label
	}

	for index, label := range result.Labels[AddressLabel] {
		m.labelsAddress.data[index] = label
	}

	return nil
}
