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
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"errors"
	"strings"
	"sync"
)

type RPCMap struct {
	mu   sync.RWMutex
	data map[uint32]*RPC
}

func (m *RPCMap) TitleExists(title string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for idx := range m.data {
		if strings.ToLower(m.data[idx].title) == strings.ToLower(title) {
			return true
		}
	}
	return false
}

func (m *RPCMap) Data() map[uint32]*RPC {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

func (m *RPCMap) Get(index uint32) *RPC {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data[index]
}

func (m *RPCMap) All() map[uint32]*RPC {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data
}

func (m *RPCMap) Add(title, endpoint string) (uint32, error) {
	var lastIndex uint32

	if m.TitleExists(title) {
		return 0, errors.New("label already exist")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for lastIndex = range m.data {
	}

	lastIndex++
	m.data[lastIndex] = &RPC{
		title:    title,
		endpoint: endpoint,
	}
	return lastIndex, nil
}

func (m *RPCMap) Remove(index uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[index]; !ok {
		return errors.New("label not exist")
	}
	delete(m.data, index)
	return nil
}
