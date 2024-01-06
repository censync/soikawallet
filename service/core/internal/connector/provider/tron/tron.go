package tron

import (
	"github.com/censync/soikawallet/service/core/internal/connector/client"
	"github.com/censync/soikawallet/service/core/internal/connector/client/evm"
	"github.com/censync/soikawallet/service/core/internal/types"
)

const (
	ProviderTypeTron = "tron"
)

type Tron struct {
	ctx    *types.RPCContext
	client *evm.ClientEVM
}

func NewTron() *Tron {
	return &Tron{}
}

func (t *Tron) GetType() string {
	return ProviderTypeTron
}

func (t *Tron) WithClient(ctx *types.RPCContext, tronClient client.Client) (*Tron, error) {
	// panic for safe
	t.client = tronClient.(*evm.ClientEVM)
	t.ctx = ctx
	return t, nil
}

func (t *Tron) GetHeight() (uint64, error) {
	return 0, nil
}
