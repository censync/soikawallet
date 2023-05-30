package dto

type GetTokensByNetworkDTO struct {
	CoinType uint32
}

type GetTokenDTO struct {
	Standard uint8
	CoinType uint32
	Contract string
}

type AddTokenDTO struct {
	Standard       uint8
	CoinType       uint32
	Contract       string
	DerivationPath string // ?null
}

type GetChainsDTO struct {
	OnlyW3 bool
}
