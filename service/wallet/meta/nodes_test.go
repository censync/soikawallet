package meta

import (
	"encoding/json"
	"github.com/censync/soikawallet/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
)

func init() {
	metaNodes = &nodes{}
	metaNodes.initNodes()
}

func TestNodes_AddRPCNode_Positive(t *testing.T) {
	// test init
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.accountsLinks)

	for index := range testDataRPC {
		rpc := types.NewRPC(testDataRPC[index][0], testDataRPC[index][1], false)

		nodeIndex := types.NodeIndex{
			CoinType: types.Ethereum,
			Index:    uint32(index + 1),
		}
		metaNodes.AddRPCNode(nodeIndex, rpc)
	}

	if len(metaNodes.nodes) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	if len(metaNodes.accountsLinks) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	for nodeIndex, rpc := range metaNodes.nodes {
		index := nodeIndex.Index - 1
		assert.Equal(t, testDataRPC[index][0], rpc.Title())
		assert.Equal(t, testDataRPC[index][1], rpc.Endpoint())
	}
}

func TestNodes_SetRPCAccountLink_Positive(t *testing.T) {
	// test init
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.accountsLinks)

	for accountIndex := types.AccountIndex(0); accountIndex < types.AccountIndex(len(testDataRPC)); accountIndex++ {
		nodeIndex := types.NodeIndex{
			CoinType: types.Ethereum,
			Index:    uint32(accountIndex + 1),
		}

		err := metaNodes.SetRPCAddressLink(nodeIndex, accountIndex)

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
	// test init
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.accountsLinks)

	for accountIndex := types.AccountIndex(0); accountIndex < types.AccountIndex(len(testDataRPC)); accountIndex++ {
		nodeIndex := types.NodeIndex{
			CoinType: types.Ethereum,
			Index:    uint32(accountIndex + 1),
		}
		exists := metaNodes.IsRPCAccountLinkExists(nodeIndex, accountIndex)

		assert.Equal(t, true, exists)

		metaNodes.RemoveRPCAccountLink(nodeIndex, accountIndex)

		notExists := metaNodes.IsRPCAccountLinkExists(nodeIndex, accountIndex)

		assert.Equal(t, false, notExists)
	}
}

func TestNodes_RemoveRPCNode_Positive(t *testing.T) {
	// test init
	assert.NotNil(t, metaNodes)
	assert.NotNil(t, metaNodes.nodes)
	assert.NotNil(t, metaNodes.accountsLinks)

	if len(metaNodes.nodes) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	if len(metaNodes.accountsLinks) != len(testDataRPC) {
		t.Fatal("incorrect length")
	}

	for index := len(testDataRPC); index > 0; index-- {
		nodeIndex := types.NodeIndex{
			CoinType: types.Ethereum,
			Index:    uint32(index),
		}

		err := metaNodes.RemoveRPCNode(nodeIndex)
		assert.Nil(t, err)

		_, exists := metaNodes.nodes[nodeIndex]
		assert.False(t, exists)

		_, exists = metaNodes.accountsLinks[nodeIndex]
		assert.False(t, exists)
	}

	if len(metaNodes.nodes) > 0 {
		t.Fatal("awaiting zero length")
	}

	if len(metaNodes.accountsLinks) > 0 {
		t.Fatal("awaiting zero length")
	}
}
