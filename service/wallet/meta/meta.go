package meta

import (
	"encoding/json"
	"errors"
	"github.com/censync/soikawallet/types"
)

const (
	metaSettingsVersion = 1
)

type Meta struct {
	version             uint8
	labelsAccount       Labels
	labelsAddress       Labels
	nodesAccountsLinks  map[types.NodeIndex][]types.AccountIndex
	labelsAccountsLinks map[uint32]types.AccountIndex
	// addressLabels map[NodeIndex][]types.AddressIndex
}

func InitMeta() *Meta {
	return &Meta{
		version:            metaSettingsVersion,
		labelsAccount:      initLabels(),
		labelsAddress:      initLabels(),
		nodesAccountsLinks: map[types.NodeIndex][]types.AccountIndex{},
	}
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Version uint8                           `json:"v"`
		Labels  map[LabelType]map[uint32]string `json:"labels"`
		// Nodes   map[types.NodeIndex]map[string]map[string]string
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

// Labels
func (m *Meta) AccountLabels() map[uint32]string {
	return m.labelsAccount.Data()
}

func (m *Meta) AddressLabels() map[uint32]string {
	return m.labelsAddress.Data()
}

func (m *Meta) AddAccountLabel(label string) (uint32, error) {
	return m.labelsAccount.Add(label)
}

func (m *Meta) AddAddressLabel(label string) (uint32, error) {
	return m.labelsAddress.Add(label)
}

func (m *Meta) RemoveAccountLabel(index uint32) error {
	return m.labelsAccount.Remove(index)
}

func (m *Meta) RemoveAddressLabel(index uint32) error {
	return m.labelsAddress.Remove(index)
}

// Linked Labels

// Linked nodes

func (m *Meta) GetRPCAccountLinks(nodeKey types.NodeIndex) []types.AccountIndex {
	return m.nodesAccountsLinks[nodeKey]
}

func (m *Meta) GetRPCAccountLinksCount(coinType types.CoinType, nodeIndex uint32) int {
	return len(m.nodesAccountsLinks[types.NodeIndex{
		CoinType: coinType,
		Index:    nodeIndex,
	}])
}

func (m *Meta) IsRPCAccountLinkExists(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) bool {
	accounts := m.GetRPCAccountLinks(nodeIndex)

	if accounts != nil {
		for _, index := range accounts {
			if index == accountIndex {
				return true
			}
		}
	}
	return false
}

func (m *Meta) SetRPCAccountLink(coinType types.CoinType, accountIndex types.AccountIndex, nodeIndex uint32) error {
	nodeKey := types.NodeIndex{
		CoinType: coinType,
		Index:    nodeIndex,
	}

	if m.IsRPCAccountLinkExists(nodeKey, accountIndex) {
		return errors.New("already enabled")
	}
	m.nodesAccountsLinks[nodeKey] = append(m.nodesAccountsLinks[nodeKey], accountIndex)
	return nil
}

func (m *Meta) RemoveRPCAccountLink(nodeKey types.NodeIndex, accountIndex types.AccountIndex) {
	for index := range m.nodesAccountsLinks[nodeKey] {
		if m.nodesAccountsLinks[nodeKey][index] == accountIndex {
			m.nodesAccountsLinks[nodeKey] = append(m.nodesAccountsLinks[nodeKey][:index], m.nodesAccountsLinks[nodeKey][index+1:]...)
		}
	}
}
