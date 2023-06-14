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
	links  map[string]uint32
}

func initLabels() label {
	return label{
		labels: map[uint32]string{},
		links:  map[string]uint32{},
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

func (m *label) IsIndexExists(index uint32) bool {
	_, ok := m.labels[index]
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

func (m *label) Remove(index uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.IsIndexExists(index) {
		return errors.New("label not exist")
	}
	delete(m.labels, index)
	return nil
}

func (m *label) GetLabel(path string) string {
	if index, ok := m.links[path]; ok {
		return m.labels[index]
	}
	return ""
}

func (m *label) SetLink(path string, index uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.IsIndexExists(index) {
		return errors.New("label not exist")
	}

	if currentIndex, ok := m.links[path]; ok && currentIndex == index {
		return errors.New("already linked")
	}

	m.links[path] = index

	return nil
}

func (m *label) RemoveLink(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.links[path]; !ok {
		return errors.New("not linked")
	}

	delete(m.links, path)

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

func (l *labels) AddressLabels() map[uint32]string {
	return l.labelsAddress.labels
}

func (l *labels) GetAccountLabel(path string) string {
	return l.labelsAccount.GetLabel(path)
}

func (l *labels) SetAccountLabelLink(path string, index uint32) error {
	return l.labelsAccount.SetLink(path, index)
}

func (l *labels) RemoveAccountLabelLink(path string) error {
	return l.labelsAccount.RemoveLink(path)
}

// Address associated labels operations

func (l *labels) GetAddressLabel(path string) string {
	return l.labelsAddress.GetLabel(path)
}

func (l *labels) AddAddressLabel(label string) (uint32, error) {
	return l.labelsAddress.Add(label)
}

func (l *labels) RemoveAddressLabel(index uint32) error {
	return l.labelsAddress.Remove(index)
}

func (l *labels) SetAddressLabelLink(path string, index uint32) error {
	return l.labelsAddress.SetLink(path, index)
}

func (l *labels) RemoveAddressLabelLink(path string) error {
	return l.labelsAddress.RemoveLink(path)
}

func (l *labels) MarshalJSON() ([]byte, error) {
	result := map[uint8]interface{}{
		types.AccountLabel: l.labelsAccount.Labels(),
		types.AddressLabel: l.labelsAddress.Labels(),
	}
	return json.Marshal(result)
}
