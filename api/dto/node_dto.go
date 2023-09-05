package dto

import mhda "github.com/censync/go-mhda"

type GetTokensByNetworkDTO struct {
	ChainKey mhda.ChainKey
}

type GetTokenDTO struct {
	Standard uint8
	ChainKey mhda.ChainKey
	Contract string
}

type AddTokenDTO struct {
	Standard uint8
	ChainKey mhda.ChainKey
	Contract string
	MhdaPath string // ?null
}

type GetChainsDTO struct {
	OnlyW3 bool
}
