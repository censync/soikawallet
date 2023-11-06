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
	"strings"
	"sync"
)

const (
	AccountLabel = 1
	AddressLabel = 2
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

func (l *label) IsLabelExists(label string) bool {
	for idx := range l.labels {
		if strings.ToLower(l.labels[idx]) == strings.ToLower(label) {
			return true
		}
	}
	return false
}

func (l *label) IsIndexExists(labelIndex uint32) bool {
	_, ok := l.labels[labelIndex]
	return ok
}

func (l *label) Labels() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return map[string]interface{}{
		"d": l.labels,
		"l": l.links,
	}
}

func (l *label) Add(label string) (uint32, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var lastIndex uint32

	if l.IsLabelExists(label) {
		return 0, errors.New("label already exist")
	}

	for lastIndex = range l.labels {
	}

	lastIndex++
	l.labels[lastIndex] = label
	return lastIndex, nil
}

func (l *label) Remove(labelIndex uint32) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.IsIndexExists(labelIndex) {
		return errors.New("label not exist")
	}
	delete(l.labels, labelIndex)

	for path, linkedIndex := range l.links {
		if linkedIndex == labelIndex {
			delete(l.links, path)
		}
	}

	return nil
}

func (l *label) GetLabel(addrIdx aIndex) string {
	if index, ok := l.links[addrIdx]; ok {
		return l.labels[index]
	}
	return ""
}

func (l *label) SetLink(addrIdx aIndex, labelIndex uint32) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.IsIndexExists(labelIndex) {
		return errors.New("label not exist")
	}

	if currentIndex, ok := l.links[addrIdx]; ok && currentIndex == labelIndex {
		return errors.New("already linked")
	}

	l.links[addrIdx] = labelIndex

	return nil
}

func (l *label) RemoveLink(addrIdx aIndex) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.links[addrIdx]; !ok {
		return errors.New("not linked")
	}

	delete(l.links, addrIdx)

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
func (m *Meta) GetAccountLabel(chainKey mhda.ChainKey, accountIndex mhda.AccountIndex) string {
	return m.labelsAccount.GetLabel(aIndex(accountIndex))
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

// TODO: Check addrKey
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
		AccountLabel: l.labelsAccount.Labels(),
		AddressLabel: l.labelsAddress.Labels(),
	}
	return json.Marshal(result)
}
