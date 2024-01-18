package provider

import (
	"context"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type Provider interface {
	GetType() string
	GetHeight(context.Context) (uint64, error)
	GetBlock(context.Context, uint64) (*ethTypes.Block, error)
}
