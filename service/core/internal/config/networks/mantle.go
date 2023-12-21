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

//go:build !testnet

package networks

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types"
)

var Mantle = types.NewNetwork(
	mhda.ETH,
	`Mantle`,
	`MNT`,
	18,
	1e9,
	"gwei",
	true,
	&types.EVMConfig{
		ChainId: 0x1388,
	},
).SetDefaultRPC(
	`https://rpc.mantle.xyz`,
	`https://explorer.mantle.xyz/`, // /block/ /address/ /tx/
)
