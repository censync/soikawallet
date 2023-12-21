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

package datafeed

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/core/internal/types/currencies"
)

// Chainlink mainnet data sources
// https://docs.chain.link/data-feeds/price-feeds/addresses/?network=ethereum&page=1#ethereum-mainnet
var evmFiat = map[currencies.DataFeedType]map[mhda.ChainKey]map[string]string{
	currencies.FeedChainLink: {
		mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1`).Key(): {
			"ETH_USD":   "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419",
			"BTC_USD":   "0xf4030086522a5beea4988f8ca5b36dbc97bee88c",
			"BNB_USD":   "0x14e613ac84a31f709eadbdf89c6cc390fdc9540a",
			"MATIC_USD": "0x7bac85a8a13a4bcd8abb3eb7d6b4d632c5a57676",
			"AVAX_USD":  "0xFF3EEb22B5E3dE6e705b44749C2559d704923FD7",
		},
		mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x89`).Key(): {
			"MATIC_USD": "0xab594600376ec9fd91f8e885dadf0ce036862de0",
			"BTC_USD":   "0xc907e116054ad103354f2d350fd2514433d57f6f",
			"ETH_USD":   "0xf9680d99d6c9589e2a93a78a04a279e509205945",
			"BNB_USD":   "0x82a6c4af830caa6c97bb504425f6a66165c2c26e",
			"AVAX_USD":  "0xe01eA2fbd8D76ee323FbEd03eB9a8625EC981A10",
		},
		mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x38`).Key(): {
			"BNB_USD":   "0x0567f2323251f0aab15c8dfb1967e4e8a7d42aee",
			"BTC_USD":   "0x264990fbd0a4796a3e3d8e37c4d5f87a3aca5ebf",
			"ETH_USD":   "0x9ef1b8c0e4f7dc8bf5719ea496883dc6401d5b2e",
			"MATIC_USD": "0x7ca57b0ca6367191c94c8914d7df09a57655905f",
			"AVAX_USD":  "0x5974855ce31EE8E1fff2e76591CbF83D7110F151",
		},
	},
}
