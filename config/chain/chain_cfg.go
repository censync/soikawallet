//go:build !testnet

package chain

import mhda "github.com/censync/go-mhda"

var (
	BitcoinChain = mhda.NewChain(mhda.Bitcoin, mhda.BTC, `bitcoin`)
	TronChain    = mhda.NewChain(mhda.TronVM, mhda.TRX, `mainnet`)
	// L1
	BinanceSmartChain = mhda.NewChain(mhda.EthereumVM, mhda.BSC, `0x38`)
	EthereumChain     = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x1`)
	PolygonChain      = mhda.NewChain(mhda.EthereumVM, mhda.MATIC, `0x89`)
	Moonbeam          = mhda.NewChain(mhda.EthereumVM, mhda.GLMR, `0x504`)

	// L2
	OptimismChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa`)
	ArbitrumChain   = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa4b1`)
	AvalancheCChain = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0xa86a`)
	BaseChain       = mhda.NewChain(mhda.EthereumVM, mhda.ETH, `0x2105`)
)
