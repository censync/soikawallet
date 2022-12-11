package types

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
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
	name      string
	symbol    string
	contract  string
	decimals  int
	isBuiltin bool
}

func (t TokenConfig) Name() string {
	return t.name
}

func (t TokenConfig) Symbol() string {
	return t.symbol
}

func (t TokenConfig) Contract() string {
	return t.contract
}

func (t TokenConfig) Decimals() int {
	return t.decimals
}

func (t TokenConfig) IsBuiltin() bool {
	return t.isBuiltin
}

func NewTokenConfig(name string, symbol string, contract string, decimals int) *TokenConfig {
	return &TokenConfig{name: name, symbol: symbol, contract: contract, decimals: decimals, isBuiltin: false}
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
	AllTokens() map[string]*TokenConfig
}

type NetworkAdapter interface {
	RPCAdapter
	//
	Address(pub *ecdsa.PublicKey) string

	// node methods
	GetBalance(ctx *RPCContext) (float64, error)
	GetERC20TokenBalance(ctx *RPCContext, contract string, decimals int) (*big.Float, error)
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

func (n *BaseNetwork) SetBuiltinToken(name, symbol, contract string, decimals int) *BaseNetwork {
	if _, ok := n.tokens[contract]; ok {
		panic(fmt.Sprintf("token \"%s\" contract \"%s\" already exist", symbol, contract))
	}
	n.tokens[contract] = &TokenConfig{
		name:      name,
		symbol:    symbol,
		contract:  contract,
		decimals:  decimals,
		isBuiltin: true,
	}
	return n
}

func (n *BaseNetwork) AddToken(symbol, contract string) error {
	if _, ok := n.tokens[contract]; ok {
		return errors.New(fmt.Sprintf("token \"%s\" contract \"%s\" already exist", symbol, contract))
	}
	n.tokens[contract] = &TokenConfig{
		symbol:    symbol,
		contract:  contract,
		isBuiltin: false,
	}
	return nil
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

func (n *BaseNetwork) AllTokens() map[string]*TokenConfig {
	return n.tokens
}
