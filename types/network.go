package types

import (
	"context"
	"crypto/ecdsa"
)

const (
	defaultRPCIndex = 0
	defaultRPCTitle = `default`
)

type BaseNetwork struct {
	// Slip 44 index
	index    CoinType
	name     string
	currency string
	explorer string

	rpc    *RPCMap
	tokens map[string]*TokenConfig

	NetworkAdapter
}

/*type EVMConfig struct {
	ChainId uint32
}*/

type RPCContext struct {
	context.Context
	coinType       CoinType
	nodeId         uint32
	currentAccount string
}

func NewRPCContext(coinType CoinType, nodeId uint32, address ...string) *RPCContext {
	var currentAccount string
	if address != nil {
		if len(address) != 1 {
			panic("to many current addresses set")
		}
		currentAccount = address[0]
	}
	return &RPCContext{
		Context:        context.Background(),
		coinType:       coinType,
		nodeId:         nodeId,
		currentAccount: currentAccount,
	}
}

func (c *RPCContext) CoinType() CoinType {
	return c.coinType
}

func (c *RPCContext) NodeId() uint32 {
	return c.nodeId
}

func (c *RPCContext) CurrentAccount() string {
	return c.currentAccount
}

type RPC struct {
	title     string
	endpoint  string
	isDefault bool
}

func (r *RPC) Title() string {
	return r.title
}

func (r *RPC) Endpoint() string {
	return r.endpoint
}

func (r *RPC) IsDefault() bool {
	return r.isDefault
}

type TokenConfig struct {
	currency  string
	contract  string
	isBuiltin bool
}

type RPCAdapter interface {
	Name() string
	Currency() string

	// RPC
	DefaultRPC() *RPC
	RPC(index uint32) *RPC
	AllRPC() map[uint32]*RPC
	AddRPC(title, endpoint string) (uint32, error)
	RemoveRPC(index uint32) error
}

type NetworkAdapter interface {
	RPCAdapter
	//
	Address(pub *ecdsa.PublicKey) string

	// node methods
	GetBalance(ctx *RPCContext) (float64, error)
	TxSendBase(ctx *RPCContext, to string, key *ecdsa.PrivateKey) (string, error)
	TxGetReceipt(ctx *RPCContext, tx string) (map[string]interface{}, error)

	// RPC info
	GetRPCInfo(ctx *RPCContext) (map[string]interface{}, error)
}

func NewNetwork(index CoinType, name, currency string) *BaseNetwork {
	return &BaseNetwork{
		index:    index,
		name:     name,
		currency: currency,
		rpc: &RPCMap{
			data: map[uint32]*RPC{},
		},
		tokens: map[string]*TokenConfig{},
	}
}

func (n *BaseNetwork) SetDefaultRPC(defaultEndpoint, defaultExplorer string) *BaseNetwork {
	n.rpc.data[defaultRPCIndex] = &RPC{
		title:     defaultRPCTitle,
		endpoint:  defaultEndpoint,
		isDefault: true,
	}
	return n
}

func (n *BaseNetwork) SetBuiltinToken(currency, contract string) *BaseNetwork {
	n.tokens[contract] = &TokenConfig{
		currency:  currency,
		contract:  contract,
		isBuiltin: true,
	}
	return n
}

func (n *BaseNetwork) Name() string {
	return n.name
}

func (n *BaseNetwork) Currency() string {
	return n.currency
}

func (n *BaseNetwork) RPC(index uint32) *RPC {
	return n.rpc.Get(index)
}

func (n *BaseNetwork) DefaultRPC() *RPC {
	return n.RPC(defaultRPCIndex)
}

func (n *BaseNetwork) AllRPC() map[uint32]*RPC {
	return n.rpc.All()
}

func (n *BaseNetwork) RemoveRPC(index uint32) error {
	return n.rpc.Remove(index)
}

func (n *BaseNetwork) AddRPC(title, endpoint string) (uint32, error) {
	return n.rpc.Add(title, endpoint)
}
