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

//go:build testnet

package chain

import mhda "github.com/censync/go-mhda"

var (
	// TODO: Provide configs with non-hardcoded config
	BitcoinChain = mhda.NewChain(mhda.Bitcoin, mhda.BTC, `bitcoin_testnet`)
	TronChain    = mhda.NewChain(mhda.TronVM, mhda.TRX, `shasta`)
	// L1
	BinanceSmartChain = mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x61`)
	EthereumChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xaa36a7`)
	PolygonChain      = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x13881`)
	Moonbeam          = mhda.NewChain(mhda.EthereumVM, mhda.GLMR, `0x507`)
	AvalancheCChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa869`)

	// L2
	// https://community.optimism.io/docs/useful-tools/networks/#op-sepolia
	OptimismChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xaa37dc`)
	//https://docs.arbitrum.io/getting-started-users
	ArbitrumChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x66eed`)
	BaseChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x14a33`)
	MantleChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1389`)
	Blast         = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa0c71fd`)

	// ZK
	ZkPolygon = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x5a2`)
	ZkSyncEra = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x118`)
	Linea     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xe704`)
)
