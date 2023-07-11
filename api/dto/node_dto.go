package dto

type GetTokensByNetworkDTO struct {
	NetworkType uint32
}

type GetTokenDTO struct {
	Standard    uint8
	NetworkType uint32
	Contract    string
}

type AddTokenDTO struct {
	Standard       uint8
	NetworkType    uint32
	Contract       string
	DerivationPath string // ?null
}

type GetChainsDTO struct {
	OnlyW3 bool
}
