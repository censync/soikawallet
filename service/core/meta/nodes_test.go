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
	json "encoding/json"
	mhda "github.com/censync/go-mhda"
	types2 "github.com/censync/soikawallet/service/core/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

const randomOffset = 200

var (
	metaNodes   *nodes
	testDataRPC = [][]string{
		{
			"Test RPC 1",
			"https://rpc1.example.com",
		},
		{
			"Test RPC 2",
			"https://rpc2.example.com",
		},
		{
			"Test RPC 3",
			"https://rpc3.example.com",
		},
		{
			"Test RPC 4",
			"https://rpc4.example.com",
		},
		{
			"Test RPC 5",
			"https://rpc5.example.com",
		},
	}
	chainKey = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1`)
)

func init() {
	metaNodes = &nodes{}
	metaNodes.initNodes()
}

func TestNodes_AddRPCNode_Positive(t *testing.T) {
	// test init_wallet
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.subIndex)
	assert.NotNil(t, metaNodes.links)

	for index := range testDataRPC {
		rpc := types2.NewRPC(testDataRPC[index][0], testDataRPC[index][1], false)

		nodeIndex := types2.NodeIndex{
			ChainKey: chainKey.Key(),
			Index:    uint32(index + 1),
		}
		err := metaNodes.AddRPCNode(nodeIndex, rpc)
		assert.Nil(t, err)
	}

	if len(metaNodes.nodes) != len(metaNodes.subIndex) {
		t.Fatal("incorrect length")
	}

	for nodeIndex, rpc := range metaNodes.nodes {
		index := nodeIndex - 1
		assert.Equal(t, testDataRPC[index][0], rpc.Title())
		assert.Equal(t, testDataRPC[index][1], rpc.Endpoint())
	}

}

func TestNodes_SetRPCAccountLink_Positive(t *testing.T) {
	// test init_wallet
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.subIndex)
	assert.NotNil(t, metaNodes.links)

	for addrIdx := aIndex(randomOffset); addrIdx < aIndex(len(testDataRPC)+randomOffset); addrIdx++ {
		nodeIndex := types2.NodeIndex{
			ChainKey: chainKey.Key(),
			Index:    uint32(addrIdx - randomOffset + 1),
		}

		err := metaNodes.SetRPCAddressLink(addrIdx, nodeIndex)

		assert.Nil(t, err)

	}
}

func TestNodes_MarshalJSON_Positive(t *testing.T) {
	// TestNodes_AddRPCNode_Positive(t)
	// TestNodes_SetRPCAccountLink_Positive(t)
	_, err := json.Marshal(metaNodes)
	assert.Nil(t, err)
}

func TestNodes_RemoveRPCAccountLink_Positive(t *testing.T) {
	// test init_wallet
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.subIndex)
	assert.NotNil(t, metaNodes.links)

	for addrIdx := aIndex(randomOffset); addrIdx < aIndex(len(testDataRPC)+randomOffset); addrIdx++ {
		nodeIndex := types2.NodeIndex{
			ChainKey: chainKey.Key(),
			Index:    uint32(addrIdx - randomOffset + 1),
		}
		exists := metaNodes.IsRPCAccountLinkExists(addrIdx, nodeIndex)

		assert.Equal(t, true, exists)

		metaNodes.RemoveRPCAccountLink(addrIdx, nodeIndex)

		notExists := metaNodes.IsRPCAccountLinkExists(addrIdx, nodeIndex)

		assert.Equal(t, false, notExists)
	}
}

func TestNodes_RemoveRPCNode_Positive(t *testing.T) {
	// test init_wallet
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.subIndex)
	assert.NotNil(t, metaNodes.links)

	if len(metaNodes.nodes) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	if len(metaNodes.subIndex) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	for index := len(testDataRPC); index > 0; index-- {
		nodeIndex := types2.NodeIndex{
			ChainKey: chainKey.Key(),
			Index:    uint32(index + 1),
		}

		err := metaNodes.RemoveRPCNode(nodeIndex)
		assert.Nil(t, err)

		//_, exists := metaNodes.nodes[nodeIndex]
		//assert.False(t, exists)

		_, exists := metaNodes.subIndex[nodeIndex]
		assert.False(t, exists)
	}

	if len(metaNodes.nodes) > 0 {
		t.Fatal("awaiting zero length")
	}

	if len(metaNodes.subIndex) > 0 {
		t.Fatal("awaiting zero length")
	}
}
