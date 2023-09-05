package meta

import (
	"encoding/json"
	"errors"
	"github.com/censync/soikawallet/types"
	"strings"
	"sync"
)

type label struct {
	mu     sync.RWMutex
	labels map[uint32]string
	links  map[aIndex]uint32
}

func initLabels() label {
	return label{
		labels: map[uint32]string{},
		links:  map[aIndex]uint32{},
	}
}

func (m *label) IsLabelExists(label string) bool {
	for idx := range m.labels {
		if strings.ToLower(m.labels[idx]) == strings.ToLower(label) {
			return true
		}
	}
	return false
}

func (m *label) IsIndexExists(labelIndex uint32) bool {
	_, ok := m.labels[labelIndex]
	return ok
}

func (m *label) Labels() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"d": m.labels,
		"l": m.links,
	}
}

func (m *label) Add(label string) (uint32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastIndex uint32

	if m.IsLabelExists(label) {
		return 0, errors.New("label already exist")
	}

	for lastIndex = range m.labels {
	}

	lastIndex++
	m.labels[lastIndex] = label
	return lastIndex, nil
}

func (m *label) Remove(labelIndex uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.IsIndexExists(labelIndex) {
		return errors.New("label not exist")
	}
	delete(m.labels, labelIndex)

	for path, linkedIndex := range m.links {
		if linkedIndex == labelIndex {
			delete(m.links, path)
		}
	}

	return nil
}

func (m *label) GetLabel(addrIdx aIndex) string {
	if index, ok := m.links[addrIdx]; ok {
		return m.labels[index]
	}
	return ""
}

func (m *label) SetLink(addrIdx aIndex, labelIndex uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.IsIndexExists(labelIndex) {
		return errors.New("label not exist")
	}

	if currentIndex, ok := m.links[addrIdx]; ok && currentIndex == labelIndex {
		return errors.New("already linked")
	}

	m.links[addrIdx] = labelIndex

	return nil
}

func (m *label) RemoveLink(addrIdx aIndex) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.links[addrIdx]; !ok {
		return errors.New("not linked")
	}

	delete(m.links, addrIdx)

	return nil
}

type labels struct {
	labelsAccount label
	labelsAddress label
}

func (l *labels) initLabels() {
	l.labelsAccount = initLabels()
	l.labelsAddress = initLabels()
}

// Account associated labels operations

func (l *labels) AccountLabels() map[uint32]string {
	return l.labelsAccount.labels
}

func (l *labels) AddAccountLabel(label string) (uint32, error) {
	return l.labelsAccount.Add(label)
}

func (l *labels) RemoveAccountLabel(index uint32) error {
	return l.labelsAccount.Remove(index)
}

// Think about accounts
func (m *Meta) GetAccountLabel(path string) string {
	//return l.labelsAccount.GetLabel(path)
	return ``
}

func (m *Meta) SetAccountLabelLink(path string, index uint32) error {
	//return l.labelsAccount.SetLink(path, index)
	return nil
}

func (m *Meta) RemoveAccountLabelLink(path string) error {
	//return l.labelsAccount.RemoveLink(path)
	return nil
}

// Address associated labels operations

func (l *labels) AddressLabels() map[uint32]string {
	return l.labelsAddress.labels
}

func (m *Meta) GetAddressLabel(addrKey string) string {
	addrOpts, ok := m.addresses[addrKey]
	if !ok {
		return ``
	}
	return m.labelsAddress.GetLabel(addrOpts.subIndex)
}

func (l *labels) AddAddressLabel(label string) (uint32, error) {
	return l.labelsAddress.Add(label)
}

func (l *labels) RemoveAddressLabel(index uint32) error {
	return l.labelsAddress.Remove(index)
}

func (m *Meta) SetAddressLabelLink(addrKey string, index uint32) error {
	addrOpts, ok := m.addresses[addrKey]
	if !ok {
		return errors.New("address not exists")
	}
	return m.labelsAddress.SetLink(addrOpts.subIndex, index)
}

func (m *Meta) RemoveAddressLabelLink(addrKey string) error {
	addrOpts, ok := m.addresses[addrKey]
	if !ok {
		return errors.New("address not exists")
	}
	return m.labelsAddress.RemoveLink(addrOpts.subIndex)
}

func (l *labels) MarshalJSON() ([]byte, error) {
	result := map[uint8]interface{}{
		types.AccountLabel: l.labelsAccount.Labels(),
		types.AddressLabel: l.labelsAddress.Labels(),
	}
	return json.Marshal(result)
}
