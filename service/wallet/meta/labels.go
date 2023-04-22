package meta

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

const (
	AccountLabel = LabelType(1)
	AddressLabel = LabelType(2)
)

type LabelType uint8

type Labels struct {
	mu   sync.RWMutex
	data map[uint32]string
}

func initLabels() Labels {
	return Labels{
		data: map[uint32]string{},
	}
}

func (m *Labels) exists(label string) bool {
	for idx := range m.data {
		if strings.ToLower(m.data[idx]) == strings.ToLower(label) {
			return true
		}
	}
	return false
}

func (m *Labels) Data() map[uint32]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

func (m *Labels) Add(label string) (uint32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastIndex uint32

	if m.exists(label) {
		return 0, errors.New("label already exist")
	}

	for lastIndex = range m.data {
	}

	lastIndex++
	m.data[lastIndex] = label
	return lastIndex, nil
}

func (m *Labels) Remove(index uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[index]; !ok {
		return errors.New("label not exist")
	}
	delete(m.data, index)
	return nil
}

type labels struct {
	labelsAccount Labels
	labelsAddress Labels
}

func (l *labels) initLabels() {
	l.labelsAccount = initLabels()
	l.labelsAddress = initLabels()
}

func (l *labels) AccountLabels() map[uint32]string {
	return l.labelsAccount.Data()
}

func (l *labels) AddressLabels() map[uint32]string {
	return l.labelsAddress.Data()
}

func (l *labels) AddAccountLabel(label string) (uint32, error) {
	return l.labelsAccount.Add(label)
}

func (l *labels) AddAddressLabel(label string) (uint32, error) {
	return l.labelsAddress.Add(label)
}

func (l *labels) RemoveAccountLabel(index uint32) error {
	return l.labelsAccount.Remove(index)
}

func (l *labels) RemoveAddressLabel(index uint32) error {
	return l.labelsAddress.Remove(index)
}

func (l *labels) MarshalJSON() ([]byte, error) {
	result := map[LabelType]interface{}{
		AccountLabel: l.labelsAccount.Data(),
		AddressLabel: l.labelsAddress.Data(),
	}
	return json.Marshal(result)
}
