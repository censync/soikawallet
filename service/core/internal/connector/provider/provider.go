package provider

import ethTypes "github.com/ethereum/go-ethereum/core/types"

type Provider interface {
	GetType() string
	GetHeight() (uint64, error)
	GetBlock(uint64) (*ethTypes.Block, error)
}
