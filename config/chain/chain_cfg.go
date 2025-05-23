// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as pu	blished by
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

package chain

import mhda "github.com/censync/go-mhda"

var (
	// TODO: Provide configs with non-hardcoded config
	BitcoinChain = mhda.NewChain(mhda.Bitcoin, mhda.BTC, `bitcoin`)
	TronChain    = mhda.NewChain(mhda.TronVM, mhda.TRX, `mainnet`)
	// L1
	BinanceSmartChain = mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x38`)
	EthereumChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1`)
	PolygonChain      = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x89`)
	Moonbeam          = mhda.NewChain(mhda.EthereumVM, mhda.GLMR, `0x504`)
	AvalancheCChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa86a`)
	Gnosis            = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x64`)

	// L2
	OptimismChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa`)
	ArbitrumChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa4b1`)
	BaseChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x2105`)
	MantleChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1388`)
	Blast         = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa0c71fd`)

	// ZK
	ZkPolygon = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x44d`)
	ZkSyncEra = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x144`)
	Linea     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xe708`)
)
