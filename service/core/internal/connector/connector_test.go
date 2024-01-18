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

package connector

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

var testChains = map[mhda.ChainKey][]string{
	"ethereum": {"wss://go.getblock.io/"},
	// "bsc":      {"wss://bsc.publicnode.com"},           // public node has sub problems "wss://bsc.publicnode.com"
	// "polygon":  {"wss://polygon.gateway.tenderly.com"}, // public node has sub problems "wss://polygon-bor.publicnode.com"
	// "moonbeam": {"wss://moonbeam.public.blastapi.io", "wss://moonbeam-rpc.dwellir.com", "wss://moonbeam.unitedbloc.com"},
	// "op":       {"wss://optimism.gateway.tenderly.co"},
	// "gnosis":   {"wss://rpc.gnosischain.com/wss"},
}

func Test_Connector(t *testing.T) {
	preparedRPCs := map[mhda.ChainKey][]*types.RPC{}

	indexRPC := uint32(1)

	for chainKey, rpcs := range testChains {
		for _, rpc := range rpcs {
			preparedRPCs[chainKey] = append(
				preparedRPCs[chainKey],
				types.NewRPC(
					indexRPC,
					fmt.Sprintf("Node %d", indexRPC),
					rpc,
					false,
				))
			indexRPC++
		}

	}
	logrus.SetLevel(logrus.InfoLevel)
	rpcConnector, err := NewConnector(preparedRPCs)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(rpcConnector)

	t.Log(rpcConnector.AvailableChains())

	time.Sleep(5 * time.Second)

	/*for chainKey := range testChains {
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(1000*time.Millisecond))
		resp, err := rpcConnector.Call(ctx, chainKey, "eth_getBlockByNumber", []interface{}{"finalized", false})
		t.Log(resp, err)
	}*/

	time.Sleep(60 * time.Second)

	// rpcConnector.Unsubscribe("ethereum", "base_height")

	//time.Sleep(120 * time.Second)
	rpcConnector.Shutdown()

	t.Log("Finished. Next messages indicates wrong processes")
	time.Sleep(5 * time.Second)
}
