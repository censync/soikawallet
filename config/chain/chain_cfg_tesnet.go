//go:build testnet

package chain

import mhda "github.com/censync/go-mhda"

var (
	BitcoinChain = mhda.NewChain(mhda.Bitcoin, mhda.BTC, `bitcoin_testnet`)
	TronChain    = mhda.NewChain(mhda.TronVM, mhda.TRX, `shasta`)
	// L1
	BinanceSmartChain = mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x61`)
	EthereumChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xaa36a7`)
	PolygonChain      = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x13881`)
	Moonbeam          = mhda.NewChain(mhda.EthereumVM, mhda.GLMR, `0x507`)

	// L2
	// https://community.optimism.io/docs/useful-tools/networks/#op-sepolia
	OptimismChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xaa37dc`)
	//https://docs.arbitrum.io/getting-started-users
	ArbitrumChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x66eed`)
	AvalancheCChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa869`)
	BaseChain       = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x14a33`)
)
