package types

const (
	OpInitBootstrap uint8 = iota + 1
	OpInitAirGap
	OpInitWallet

	OpMetaAirGap = 10
	OpMetaWallet = 11

	OpDeriveAirGap = 30
	OpDeriveWallet = 31

	OpTxSend  = 50
	OpTxSwap5 = 51
	OpTxWeb3  = 52
)
