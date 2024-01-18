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
)

func Test_Provider(t *testing.T) {
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
	defer rpcConnector.Shutdown()

	if err != nil {
		t.Fatal(err)
	}

	ctx := types.NewRPCContext("ethereum", 0)
	provider, err := rpcConnector.GetProvider(ctx)
	// lastBlock, err := provider.GetHeight()

	block, err := provider.GetBlock(ctx, 19026186)

	t.Log("GetHeight", block, err)

	//time.Sleep(3 * time.Second)

}
