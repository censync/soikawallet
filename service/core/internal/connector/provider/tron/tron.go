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

package tron

import (
	"context"
	"github.com/censync/soikawallet/service/core/internal/connector/client"
	"github.com/censync/soikawallet/service/core/internal/connector/client/evm"
	"github.com/censync/soikawallet/service/core/internal/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

const (
	ProviderTypeTron = "tron"
)

type Tron struct {
	ctx    *types.RPCContext
	client *evm.ClientEVM
}

func NewTron() *Tron {
	return &Tron{}
}

func (t *Tron) GetType() string {
	return ProviderTypeTron
}

func (t *Tron) WithClient(tronClient client.Client) (*Tron, error) {
	// panic for safe
	t.client = tronClient.(*evm.ClientEVM)
	return t, nil
}

func (t *Tron) GetHeight(ctx context.Context) (uint64, error) {
	return 0, nil
}

func (t *Tron) GetBlock(ctx context.Context, _ uint64) (*ethTypes.Block, error) {
	return nil, nil
}
