package meta

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/types"
	"sync"
)

type nodes struct {
	mu            sync.RWMutex
	nodes         map[types.NodeIndex]*types.RPC
	accountsLinks map[types.NodeIndex][]types.AccountIndex
}

func (n *nodes) initNodes() {
	n.nodes = map[types.NodeIndex]*types.RPC{}
	n.accountsLinks = map[types.NodeIndex][]types.AccountIndex{}
}

func (n *nodes) AddRPCNode(index types.NodeIndex, rpc *types.RPC) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	n.nodes[index] = rpc
	if _, ok := n.accountsLinks[index]; !ok {
		n.accountsLinks[index] = []types.AccountIndex{}
	}
}

func (n *nodes) RemoveRPCNode(index types.NodeIndex) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if _, ok := n.nodes[index]; !ok {
		return errors.New("node is not IsLabelExists")
	}

	delete(n.nodes, index)
	delete(n.accountsLinks, index)

	return nil
}

// Linked accounts

func (n *nodes) IsRPCAccountLinkExists(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) bool {
	if n.accountsLinks[nodeIndex] != nil {
		for _, index := range n.accountsLinks[nodeIndex] {
			if index == accountIndex {
				return true
			}
		}
	}
	return false
}

func (n *nodes) GetRPCAccountLinks(nodeIndex types.NodeIndex) []types.AccountIndex {
	return n.accountsLinks[nodeIndex]
}

func (n *nodes) GetRPCAccountLinksCount(networkType types.NetworkType, nodeIndex uint32) int {
	return len(n.accountsLinks[types.NodeIndex{
		NetworkType: networkType,
		Index:       nodeIndex,
	}])
}

func (n *nodes) SetRPCAccountLink(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if n.IsRPCAccountLinkExists(nodeIndex, accountIndex) {
		return errors.New("already enabled")
	}

	n.accountsLinks[nodeIndex] = append(n.accountsLinks[nodeIndex], accountIndex)
	return nil
}

func (n *nodes) RemoveRPCAccountLink(nodeIndex types.NodeIndex, accountIndex types.AccountIndex) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for index := range n.accountsLinks[nodeIndex] {
		if n.accountsLinks[nodeIndex][index] == accountIndex {
			n.accountsLinks[nodeIndex] = append(n.accountsLinks[nodeIndex][:index], n.accountsLinks[nodeIndex][index+1:]...)
		}
	}
}

func (n *nodes) MarshalJSON() ([]byte, error) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	nodesExport := map[string]*types.RPC{}
	for nodeIndex, node := range n.nodes {
		nodesExport[fmt.Sprintf("%d:%d", nodeIndex.NetworkType, nodeIndex.Index)] = node
	}

	linksExport := map[string][]types.AccountIndex{}
	for nodeIndex, link := range n.accountsLinks {
		if len(link) > 0 {
			linksExport[fmt.Sprintf("%d:%d", nodeIndex.NetworkType, nodeIndex.Index)] = link
		}
	}

	result := map[string]interface{}{
		"nodes": nodesExport,
		"links": linksExport,
	}

	return json.Marshal(result)
}
