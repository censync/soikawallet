package types

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/types/gas"
	"math/big"
)

const (
	defaultRPCIndex = 0
	defaultRPCTitle = `default`
)

type BaseNetwork struct {
	// Slip 44 index
	index    NetworkType
	name     string
	currency string
	decimals int
	explorer string

	// gas
	gasUnits  uint64
	gasSymbol string

	rpc       *RPCMap
	gasCalc   *gas.Calculator
	tokens    map[string]*TokenConfig
	evmConfig *EVMConfig
	isW3      bool
	NetworkAdapter
}

type EVMConfig struct {
	ChainId uint32
}

type RPCContext struct {
	context.Context
	network        NetworkType
	nodeId         uint32
	currentAccount string
}

func NewRPCContext(network NetworkType, nodeId uint32, address ...string) *RPCContext {
	var currentAccount string
	if address != nil {
		if len(address) != 1 {
			panic("to many current addresses set")
		}
		currentAccount = address[0]
	}
	return &RPCContext{
		Context:        context.Background(),
		network:        network,
		nodeId:         nodeId,
		currentAccount: currentAccount,
	}
}

func (c *RPCContext) Network() NetworkType {
	return c.network
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

func NewRPC(title, endpoint string, isDefault bool) *RPC {
	return &RPC{
		title:     title,
		endpoint:  endpoint,
		isDefault: isDefault,
	}
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

func (r *RPC) MarshalJSON() ([]byte, error) {
	result := fmt.Sprintf(
		`["%s","%s","%t"]`,
		r.Title(),
		r.Endpoint(),
		r.isDefault,
	)
	return []byte(result), nil
}

type TokenConfig struct {
	standard  TokenStandard
	name      string
	symbol    string
	contract  string
	decimals  int
	isBuiltin bool
}

func (t *TokenConfig) Standard() TokenStandard {
	return t.standard
}

func (t *TokenConfig) Name() string {
	return t.name
}

func (t *TokenConfig) Symbol() string {
	return t.symbol
}

func (t *TokenConfig) Contract() string {
	return t.contract
}

func (t *TokenConfig) Decimals() int {
	return t.decimals
}

func (t *TokenConfig) IsBuiltin() bool {
	return t.isBuiltin
}

func (t *TokenConfig) MarshalJSON() ([]byte, error) {
	result := fmt.Sprintf(
		`["%d","%s","%s","%s","%d","%t"]`,
		t.Standard(),
		t.Name(),
		t.Symbol(),
		t.Contract(),
		t.Decimals(),
		t.IsBuiltin(),
	)
	return []byte(result), nil
}

func NewTokenConfig(standard TokenStandard, name string, symbol string, contract string, decimals int) *TokenConfig {
	return &TokenConfig{
		standard:  standard,
		name:      name,
		symbol:    symbol,
		contract:  contract,
		decimals:  decimals,
		isBuiltin: false,
	}
}

type RPCAdapter interface {
	Name() string
	Currency() string
	Decimals() int
	IsW3() bool
	EVMConfig() *EVMConfig

	// RPC
	DefaultRPC() *RPC
	DefaultNodeId() uint32
	RPC(index uint32) *RPC
	AllRPC() map[uint32]*RPC
	AddRPC(title, endpoint string) (uint32, error)
	RemoveRPC(index uint32) error

	GetBaseToken() *TokenConfig
	GetAllTokens() map[string]*TokenConfig
	// tokens
	IsTokenConfigExists(contract string) bool
	AddTokenConfig(standard TokenStandard, name, symbol, contract string, decimals int) (*TokenConfig, error)
	GetTokenConfig(contract string) *TokenConfig
}

type NetworkAdapter interface {
	RPCAdapter
	//
	Address(pub *ecdsa.PublicKey) string

	// node methods
	GetBalance(ctx *RPCContext) (float64, error)
	GetTokenBalance(ctx *RPCContext, contract string, decimals int) (*big.Float, error)
	GetToken(ctx *RPCContext, contract string) (*TokenConfig, error)
	GetTokenAllowance(ctx *RPCContext, contract, to string) (uint64, error)
	GetGasConfig(ctx *RPCContext, args ...interface{}) (map[string]uint64, error)
	TxSendBase(ctx *RPCContext, to string, value float64, gasTipCap, gasFeeCap uint64, key *ecdsa.PrivateKey) (string, error)
	TxSendBaseLegacy(ctx *RPCContext, to string, value float64, gasPrice uint64, key *ecdsa.PrivateKey) (string, error)
	TxSendToken(ctx *RPCContext, to string, value float64, token *TokenConfig, key *ecdsa.PrivateKey) (string, error)
	TxGetReceipt(ctx *RPCContext, tx string) (map[string]interface{}, error)

	// Tmp
	TxGasUnitsApprove(ctx *RPCContext, value float64, token *TokenConfig) (uint64, error)
	TxGasUnitsTransfer(ctx *RPCContext, to string, value float64, token *TokenConfig) (uint64, error)

	// Chainlink
	ChainLinkGetPrice(ctx *RPCContext, contract string) (uint64, uint8, error)

	// RPC info
	GetRPCInfo(ctx *RPCContext) (map[string]interface{}, error)
}

func NewNetwork(index NetworkType, name, currency string, decimals int, gasUnits uint64, gasSymbol string, isW3 bool, evmConfig *EVMConfig) *BaseNetwork {
	return &BaseNetwork{
		index:     index,
		name:      name,
		currency:  currency,
		decimals:  decimals,
		gasUnits:  gasUnits,
		gasSymbol: gasSymbol,
		isW3:      isW3,
		rpc: &RPCMap{
			data: map[uint32]*RPC{},
		},
		tokens:    map[string]*TokenConfig{},
		evmConfig: evmConfig,
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

func (n *BaseNetwork) addToken(standard TokenStandard, name, symbol, contract string, decimals int, isBuiltin bool) *TokenConfig {
	n.tokens[contract] = &TokenConfig{
		standard:  standard,
		name:      name,
		symbol:    symbol,
		contract:  contract,
		decimals:  decimals,
		isBuiltin: isBuiltin,
	}
	return n.tokens[contract]
}

func (n *BaseNetwork) SetBuiltinToken(standard TokenStandard, name, symbol, contract string, decimals int) *BaseNetwork {
	if _, ok := n.tokens[contract]; ok {
		panic(fmt.Sprintf("token \"%s\" contract \"%s\" already exist", symbol, contract))
	}
	n.addToken(standard, name, symbol, contract, decimals, true)
	return n
}

func (n *BaseNetwork) AddTokenConfig(standard TokenStandard, name, symbol, contract string, decimals int) (*TokenConfig, error) {
	if _, ok := n.tokens[contract]; ok {
		return nil, errors.New(fmt.Sprintf("token \"%s\" contract \"%s\" already exist", symbol, contract))
	}
	tokenConfig := n.addToken(standard, name, symbol, contract, decimals, false)
	return tokenConfig, nil
}

func (n *BaseNetwork) IsTokenConfigExists(contract string) bool {
	_, isExists := n.tokens[contract]
	return isExists
}

func (n *BaseNetwork) Name() string {
	return n.name
}

func (n *BaseNetwork) Currency() string {
	return n.currency
}

func (n *BaseNetwork) Decimals() int {
	return n.decimals
}

func (n *BaseNetwork) GasUnits() uint64 {
	return n.gasUnits
}

func (n *BaseNetwork) GasSymbol() string {
	return n.gasSymbol
}

func (n *BaseNetwork) IsW3() bool {
	return n.isW3
}

func (n *BaseNetwork) EVMConfig() *EVMConfig {
	return n.evmConfig
}

func (n *BaseNetwork) RPC(index uint32) *RPC {
	return n.rpc.Get(index)
}

func (n *BaseNetwork) DefaultRPC() *RPC {
	return n.RPC(defaultRPCIndex)
}

func (n *BaseNetwork) DefaultNodeId() uint32 {
	return defaultRPCIndex
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

func (n *BaseNetwork) GetTokenConfig(contract string) *TokenConfig {
	return n.tokens[contract]
}

func (n *BaseNetwork) GetBaseToken() *TokenConfig {
	return &TokenConfig{
		standard:  TokenBase,
		name:      n.name,
		symbol:    n.currency,
		contract:  ContractBase,
		decimals:  n.decimals,
		isBuiltin: true,
	}
}

func (n *BaseNetwork) GetAllTokens() map[string]*TokenConfig {
	return n.tokens
}
