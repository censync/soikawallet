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

package responses

import (
	mhda "github.com/censync/go-mhda"
)

// exported type
type AddressResponse struct {
	Address        string
	Path           string
	IsExternal     bool
	ChainKey       mhda.ChainKey
	AddressIndex   mhda.AddressIndex
	Account        mhda.AccountIndex
	Label          string
	IsW3           bool
	IsKeyDelivered bool
}

type AccountResponse struct {
	// Path        string
	ChainKey mhda.ChainKey
	Account  mhda.AccountIndex
	Label    string
}

type AirGapMessage struct {
	Chunks []string
}

type AddressTokenEntry struct {
	Standard string
	Name     string
	Symbol   string
	Contract string
}

type BaseCurrency struct {
	Symbol   string
	Decimals int
}

type AddressTokensListResponse map[string]*AddressTokenEntry

type AddressTokenBalanceEntry struct {
	Token   *AddressTokenEntry
	Balance float64
}

type AddressTokensBalanceListResponse map[string]*AddressTokenBalanceEntry

type TokenConfig struct {
	Standard string
	Name     string
	Symbol   string
	Contract string
	Decimals int
}

type ChainInfo struct {
	Name     string        `json:"name"`
	ChainKey mhda.ChainKey `json:"chain_key"`
}

type CalculatorConfig struct {
	Calculator []byte
}

type RPCInfo struct {
	Title     string
	Endpoint  string
	IsDefault bool
}
