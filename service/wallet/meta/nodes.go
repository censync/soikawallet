package meta

import (
	"errors"
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
	n.accountsLinks[index] = []types.AccountIndex{}
}

func (n *nodes) RemoveRPCNode(nodeIndex types.NodeIndex) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	if _, ok := n.nodes[nodeIndex]; !ok {
		return errors.New("node is not exists")
	}

	delete(n.nodes, nodeIndex)
	delete(n.accountsLinks, nodeIndex)
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

func (n *nodes) GetRPCAccountLinksCount(coinType types.CoinType, nodeIndex uint32) int {
	return len(n.accountsLinks[types.NodeIndex{
		CoinType: coinType,
		Index:    nodeIndex,
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
