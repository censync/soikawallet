// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package api_web3

import mhda "github.com/censync/go-mhda"

type GetAccountsRequest struct {
	ChainKey mhda.ChainKey `json:"chain_key"`
}

type RPCRequest struct {
	ChainKey mhda.ChainKey `json:"chain_key"`
	Method   string        `json:"request"`
	Params   []string      `json:"params"`
}
