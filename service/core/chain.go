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

package core

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/core/internal/network"
)

func (s *Wallet) GetAllEvmW3Chains() []*resp.ChainInfo {
	var result []*resp.ChainInfo
	for chainKey, provider := range network.GetAll() {
		if provider.IsW3() {
			result = append(result, &resp.ChainInfo{
				ChainKey: chainKey,
				Name:     provider.Name(),
			})
		}
	}
	return result
}
