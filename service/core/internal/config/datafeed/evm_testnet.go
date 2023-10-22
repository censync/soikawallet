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

package datafeed

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types/currencies"
)

// Chainlink testnet data sources
// https://docs.chain.link/data-feeds/price-feeds/addresses/?network=ethereum&page=1#sepolia-testnet
var evmFiat = map[currencies.DataFeedType]map[mhda.ChainKey]map[string]string{
	currencies.FeedChainLink: {
		mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xaa36a7`).Key(): {
			"ETH_USD": "0x694AA1769357215DE4FAC081bf1f309aDC325306",
			"BTC_USD": "0x1b44F3514812d835EB1BDB0acB33d3fA3351Ee43",
		},
		mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x13881`).Key(): {
			"MATIC_USD": "0xd0D5e3DB44DE05E9F294BB0a3bEEaF030DE24Ada",
			"BTC_USD":   "0x007A22900a3B98143368Bd5906f8E17e9867581b",
			"ETH_USD":   "0x0715A7794a1dc8e42615F059dD6e406A6594651A",
		},
		mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x61`).Key(): {
			"BNB_USD":   "0x2514895c72f50D8bd4B4F9b1110F0D6bD2c97526",
			"BTC_USD":   "0x5741306c21795FdCBb9b265Ea0255F499DFe515C",
			"ETH_USD":   "0x143db3CEEfbdfe5631aDD3E50f7614B6ba708BA7",
			"MATIC_USD": "0x957Eb0316f02ba4a9De3D308742eefd44a3c1719",
		},
	},
}
