package meta

import (
	"encoding/json"
	"errors"
)

const (
	metaSettingsVersion = 1
)

type Meta struct {
	version uint8
	labels
	nodes
	tokens
	//nodesAccountsLinks map[types.NodeIndex][]types.AccountIndex
	// labelsAccountsLinks map[uint32]types.AccountIndex // TODO: check for coins index
	// addressLabels map[NodeIndex][]types.
	// tokensRegistry map[types.CoinType]map[uint32]types.TokenConfig
	// addressTokens  map[string][]uint32
}

func InitMeta() *Meta {
	instance := &Meta{
		version: metaSettingsVersion,
	}

	instance.initLabels()

	instance.initNodes()

	instance.initTokens()

	return instance
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version uint8                           `json:"v"`
		Labels  map[LabelType]map[uint32]string `json:"labels"`
		// nodes   map[types.NodeIndex]map[string]map[string]string
	}{
		Version: m.version,
		Labels: map[LabelType]map[uint32]string{
			AccountLabel: m.labelsAccount.Data(),
			AddressLabel: m.labelsAddress.Data(),
		},
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
