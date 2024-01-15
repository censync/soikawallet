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

package evm_arbitrum

import (
	"fmt"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_legacy"
	"github.com/censync/soikawallet/service/core/internal/types"
)

type EVMArbitrum struct {
	*evm_legacy.EVMLegacy
}

func NewEVMArbitrum(baseNetwork *types.BaseNetwork) *EVMArbitrum {
	return &EVMArbitrum{EVMLegacy: evm_legacy.NewEVMLegacy(baseNetwork)}
}

func (e *EVMArbitrum) GetGasConfig(ctx *types.RPCContext, args ...interface{}) (map[string]uint64, error) {
	result := map[string]uint64{
		"base_fee":     0,
		"priority_fee": 0,
		"units":        7000000,
		"gas_limit":    0,
		"gas_used":     0,
	}

	gasTipCap, err := e.GetGasTipCap(ctx) // Fix to client.EstimateGas * client.SuggestGasPrice
	if err != nil {
		return result, err
	}
	if gasTipCap != nil {
		result["priority_fee"] = gasTipCap.Uint64()
	}

	gasPrice, err := e.GetGasPrice(ctx)

	if err != nil {
		return nil, err
	}

	// gas_price
	result["base_fee"] = gasPrice.Uint64()

	if len(args) > 0 {
		switch args[0].(string) {
		case "approve":
			result["units"], err = e.GasEstimateApprove(ctx, args[1].(string), args[2].(string), args[3].(*types.TokenConfig))
		case "transfer":
			result["units"], err = e.GasEstimateTransfer(ctx, args[1].(string), args[2].(string), args[3].(*types.TokenConfig))
		default:
			return nil, fmt.Errorf("wrong methond: %s", args[0])
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
