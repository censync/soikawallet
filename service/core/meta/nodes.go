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
	types2 "github.com/censync/soikawallet/service/core/internal/types"
	"sync"
)

type nodes struct {
	mu       sync.RWMutex
	nodes    map[uint32]*types2.RPC
	subIndex map[types2.NodeIndex]uint32
	links    map[aIndex][]uint32
}

func (n *nodes) initNodes() {
	n.nodes = map[uint32]*types2.RPC{}
	n.subIndex = map[types2.NodeIndex]uint32{}
	n.links = map[aIndex][]uint32{}

	// n.accountsLinks = map[types.NodeIndex][]types.AccountIndex{}
}

func (n *nodes) AddRPCNode(index types2.NodeIndex, rpc *types2.RPC) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if _, ok := n.subIndex[index]; ok {
		return errors.New("already exist")
	}

	var lastIndex uint32

	// Take max value of internal index
	for _, lastIndex = range n.subIndex {
	}

	lastIndex++

	n.subIndex[index] = lastIndex
	n.nodes[lastIndex] = rpc

	return nil
}

func (n *nodes) RemoveRPCNode(nodeIndex types2.NodeIndex) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return errors.New("node is not exists")
	}

	delete(n.nodes, internalIndex)
	delete(n.subIndex, nodeIndex)
	// TODO: Add check links

	return nil
}

// Linked accounts

func (n *nodes) IsRPCAccountLinkExists(addrIdx aIndex, nodeIndex types2.NodeIndex) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return false
	}

	if _, ok = n.links[addrIdx]; ok {
		for index := range n.links[addrIdx] {
			if n.links[addrIdx][index] == internalIndex {
				return true
			}
		}
	}
	return false
}

// WTF?
func (n *nodes) GetRPCAccountLinks(nodeIndex types2.NodeIndex) []aIndex {
	n.mu.RLock()
	defer n.mu.RUnlock()

	var result []aIndex

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return result
	}

	if len(n.links) > 0 {
		for index := range n.links {
			if len(n.links[index]) > 0 {
				for addrIdx := range n.links {
					// TODO: Add mutex
					if n.links[index][addrIdx] == internalIndex {
						result = append(result, addrIdx)
					}
				}
			}
		}
	}
	return result
}

func (n *nodes) GetRPCAccountLinksCount(nodeIndex types2.NodeIndex) int {
	var result int

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return 0
	}

	// Not tested
	for addrIdx := range n.links {
		for idx := range n.links[addrIdx] {
			if n.links[addrIdx][idx] == internalIndex {
				result++
			}
		}
	}

	return result
}

func (n *nodes) SetRPCAddressLink(addrIdx aIndex, nodeIndex types2.NodeIndex) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return errors.New("rpc not exists")
	}

	if n.IsRPCAccountLinkExists(addrIdx, nodeIndex) {
		return errors.New("already enabled")
	}

	n.links[addrIdx] = append(n.links[addrIdx], internalIndex)
	return nil
}

func (n *nodes) RemoveRPCAccountLink(addrIdx aIndex, nodeIndex types2.NodeIndex) error {
	var linkExists bool

	n.mu.Lock()
	defer n.mu.Unlock()

	internalIndex, ok := n.subIndex[nodeIndex]

	if !ok {
		return errors.New("rpc not exists")
	}

	if _, ok = n.links[addrIdx]; ok {
		for index := range n.links[addrIdx] {
			if n.links[addrIdx][index] == internalIndex {
				linkExists = true
				n.links[addrIdx] = append(n.links[addrIdx][:index], n.links[addrIdx][index+1:]...)
			}
		}
	}

	if linkExists {
		return errors.New("link not exists")
	}

	return nil
}

func (n *nodes) MarshalJSON() ([]byte, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	nodesExport := map[string]*types2.RPC{}
	/*for nodeIndex, node := range n.nodes {
		nodesExport[fmt.Sprintf("%d:%d", nodeIndex.CoinType, nodeIndex.Index)] = node
	}

	linksExport := map[string][]types.AccountIndex{}
	for nodeIndex, link := range n.accountsLinks {
		if len(link) > 0 {
			linksExport[fmt.Sprintf("%d:%d", nodeIndex.CoinType, nodeIndex.Index)] = link
		}
	}*/

	result := map[string]interface{}{
		"nodes": nodesExport,
		// "links": linksExport,
	}

	return json.Marshal(result)
}
