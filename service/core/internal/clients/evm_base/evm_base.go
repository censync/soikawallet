// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

package evm_base

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/service/core/internal/abi/oracle/chainlink"
	"github.com/censync/soikawallet/service/core/internal/abi/tokens/erc20"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
	"strings"
)

const (
	wei         = uint64(1e18)
	gwei        = uint64(1e9)
	gasMinLimit = 21000
)

var abiMap = map[string]string{
	"approve":  `[{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`,
	"transfer": `[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`,
}

type EVM struct {
	*types.BaseNetwork
	clients map[uint32]*ethclient.Client
}

func NewEVM(baseNetwork *types.BaseNetwork) *EVM {
	return &EVM{BaseNetwork: baseNetwork, clients: map[uint32]*ethclient.Client{}}
}

func (e *EVM) GetClient(nodeId uint32) (*ethclient.Client, error) {
	var err error
	if e.clients[nodeId] != nil {
		return e.clients[nodeId], nil
	} else {
		e.clients[nodeId], err = ethclient.Dial(e.DefaultRPC().Endpoint())
		return e.clients[nodeId], err
	}
}

func (e *EVM) Address(pub *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*pub).Hex()
}

func (e *EVM) GetHeight(ctx *types.RPCContext) (uint64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	return client.BlockNumber(ctx)
}

func (e *EVM) GetBlock(ctx *types.RPCContext, blockNumber uint64) (*ethTypes.Block, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}
	return client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

func (e *EVM) GetBalance(ctx *types.RPCContext) (float64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	balance, err := client.BalanceAt(ctx, common.HexToAddress(ctx.CurrentAccount()), nil)
	if err != nil {
		return 0, err
	}
	return float64(balance.Uint64()) / float64(wei), nil
}

func (e *EVM) GetTokenBalance(ctx *types.RPCContext, contract string, decimals int) (*big.Float, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}
	instance, err := erc20.NewErc20(common.HexToAddress(contract), client)
	if err != nil {
		return nil, err
	}

	balance, err := instance.BalanceOf(&bind.CallOpts{}, common.HexToAddress(ctx.CurrentAccount()))

	if err != nil {
		return nil, err
	}

	floatBalance := new(big.Float)
	floatBalance.SetString(balance.String())
	humanBalance := new(big.Float).Quo(floatBalance, big.NewFloat(math.Pow10(decimals)))

	return humanBalance, nil
}

func (e *EVM) GetToken(ctx *types.RPCContext, contract string) (*types.TokenConfig, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}
	instance, err := erc20.NewErc20(common.HexToAddress(contract), client)
	if err != nil {
		return nil, err
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})

	if err != nil {
		return nil, err
	}

	return types.NewTokenConfig(types.TokenERC20, name, symbol, contract, int(decimals)), nil
}

func (e *EVM) GetTokenAllowance(ctx *types.RPCContext, contract, to string) (uint64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	instance, err := erc20.NewErc20(common.HexToAddress(contract), client)
	if err != nil {
		return 0, err
	}

	allowance, err := instance.Allowance(nil, common.HexToAddress(ctx.CurrentAccount()), common.HexToAddress(to))

	if err != nil {
		return 0, err
	}

	return allowance.Uint64(), nil
}

func (e *EVM) GetGasPrice(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.SuggestGasPrice(ctx)
}

func (e *EVM) GetPendingNonce(ctx *types.RPCContext) (uint64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	// Known issue: Optimisms pending transactions returns error "replacement transaction underpriced"
	return client.PendingNonceAt(ctx, common.HexToAddress(ctx.CurrentAccount()))
}

func (e *EVM) GetChainId(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.ChainID(ctx)
}

func (e *EVM) GetGasTipCap(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.SuggestGasTipCap(ctx)
}

// Gas operations

func (e *EVM) GetGasEstimate(ctx *types.RPCContext, msg *ethereum.CallMsg) (uint64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}

	return client.EstimateGas(ctx, *msg)
}

func (e *EVM) GetGasConfig(ctx *types.RPCContext, args ...interface{}) (map[string]uint64, error) {
	result := map[string]uint64{
		"base_fee":     0,
		"priority_fee": 0,
		"units":        21000,
		"gas_limit":    0,
		"gas_used":     0,
	}

	height, err := e.GetHeight(ctx)

	if err != nil {
		return result, err
	}

	// Not working for L2

	block, err := e.GetBlock(ctx, height)

	if err != nil {
		return result, err
	}

	result["gas_used"] = block.GasUsed()

	gasLimit := block.GasLimit()
	result["gas_limit"] = gasLimit

	baseFee := block.BaseFee()
	if baseFee != nil {
		result["base_fee"] = baseFee.Uint64()
	}

	gasTipCap, err := e.GetGasTipCap(ctx)
	if err != nil {
		return result, err
	}
	if gasTipCap != nil {
		result["priority_fee"] = gasTipCap.Uint64()
	}

	/*
		gasPrice, err := e.GetGasPrice(ctx)

		if err != nil {
			return nil, err
		}

		// gas_price
		result["base_fee"] = gasPrice.Uint64()
	*/

	if len(args) > 0 {
		switch args[0].(string) {
		case "approve":
			result["units"], err = e.GasEstimateApprove(ctx, args[1].(string), args[2].(string), args[3].(*types.TokenConfig))
		case "transfer":
			result["units"], err = e.GasEstimateTransfer(ctx, args[1].(string), args[2].(string), args[3].(*types.TokenConfig))
		default:
			return nil, fmt.Errorf("wrong methond: %s", args[0])
		}
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func GasCalcPrepared(methodName string, args ...interface{}) []byte {
	contractABI, ok := abiMap[methodName]

	if !ok {
		panic("method not defined in prepared ABI")
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))

	if err != nil {
		panic(err)
	}

	method, ok := parsedABI.Methods[methodName]

	if !ok {
		panic("method not found")
	}

	data, err := method.Inputs.Pack(args...)

	if err != nil {
		panic(err)
	}
	methodID := crypto.Keccak256Hash([]byte(method.Sig)).Bytes()[:4]

	callData := append(methodID, data...)

	return callData
}

func (e *EVM) GasEstimateApprove(ctx *types.RPCContext, spender, value string, token *types.TokenConfig) (uint64, error) {
	addrSpender := common.HexToAddress(spender)
	weiAmount, err := StrToWei(value)

	if err != nil {
		return 0, err
	}

	callData := GasCalcPrepared("approve", addrSpender, weiAmount)

	tokenContract := common.HexToAddress(token.Contract())

	addrFrom := common.HexToAddress(ctx.CurrentAccount())

	dataTx := ethereum.CallMsg{
		From: addrFrom,
		To:   &tokenContract,
		Data: callData,
	}

	gas, err := e.GetGasEstimate(ctx, &dataTx)
	return gas, err
}

func (e *EVM) GasEstimateTransfer(ctx *types.RPCContext, to, value string, token *types.TokenConfig) (uint64, error) {
	addrTo := common.HexToAddress(to)
	weiAmount, err := StrToWei(value)

	if err != nil {
		return 0, err
	}

	callData := GasCalcPrepared("transfer", addrTo, weiAmount)

	tokenContract := common.HexToAddress(token.Contract())

	addrFrom := common.HexToAddress(ctx.CurrentAccount())

	dataTx := ethereum.CallMsg{
		From: addrFrom,
		To:   &tokenContract,
		Data: callData,
	}

	gas, err := e.GetGasEstimate(ctx, &dataTx)
	return gas, err
}

// Transactions

func (e *EVM) TxSendBase(ctx *types.RPCContext, to string, value string, gas, gasTipCap, gasFeeCap uint64, key *ecdsa.PrivateKey) (interface{}, error) {
	var txData ethTypes.TxData
	chainId, err := e.GetChainId(ctx)

	if err != nil {
		return "", err
	}

	nonce, err := e.GetPendingNonce(ctx)

	if err != nil {
		return "", err
	}

	addrTo := common.HexToAddress(to)
	weiValue, err := StrToWei(value)

	if err != nil {
		return 0, err
	}

	txData = &ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: new(big.Int).SetUint64(gasTipCap), // gasTipCap = (priorityFee)  maxPriorityFeePerGas
		GasFeeCap: new(big.Int).SetUint64(gasFeeCap), // a.k.a. maxFeePerGas limit commission gasFeeCap = gasTipCap + pendingHeader.BaseFee * 2
		Gas:       gas,                               // units
		Nonce:     nonce,
		To:        &addrTo,
		Value:     weiValue,
		Data:      nil,
	}

	tx := ethTypes.NewTx(txData)

	// AirGap
	if key == nil {
		return tx.Data(), nil
	}

	signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

	if err != nil {
		return ``, nil
	}

	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, signedTX)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) TxSendToken(ctx *types.RPCContext, to, value string, token *types.TokenConfig, gas, gasTipCap, gasFeeCap uint64, key *ecdsa.PrivateKey) (interface{}, error) {
	var txData ethTypes.TxData

	chainId, err := e.GetChainId(ctx)

	if err != nil {
		return "", err
	}

	nonce, err := e.GetPendingNonce(ctx)

	if err != nil {
		return "", err
	}

	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}

	addrTo := common.HexToAddress(to)
	weiAmount, err := StrToWei(value)

	if err != nil {
		return 0, err
	}

	callData := GasCalcPrepared("transfer", addrTo, weiAmount)

	tokenContract := common.HexToAddress(token.Contract())

	txData = &ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: new(big.Int).SetUint64(gasTipCap), // gasTipCap = (priorityFee)  maxPriorityFeePerGas
		GasFeeCap: new(big.Int).SetUint64(gasFeeCap),
		Gas:       gas,
		Nonce:     nonce,
		To:        &tokenContract,
		Data:      callData,
	}

	tx := ethTypes.NewTx(txData)

	// AirGap
	if key == nil {
		return tx.Data(), nil
	}

	signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, signedTX)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) TxApproveToken(ctx *types.RPCContext, spender string, value string, token *types.TokenConfig, gas, gasTipCap, gasFeeCap uint64, key *ecdsa.PrivateKey) (interface{}, error) {
	chainId, err := e.GetChainId(ctx)

	if err != nil {
		return "", err
	}

	nonce, err := e.GetPendingNonce(ctx)

	if err != nil {
		return "", err
	}

	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}

	addrSpender := common.HexToAddress(spender)
	weiAmount, err := StrToWei(value)

	if err != nil {
		return 0, err
	}

	callData := GasCalcPrepared("approve", addrSpender, weiAmount)

	tokenContract := common.HexToAddress(token.Contract())

	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasTipCap: new(big.Int).SetUint64(gasTipCap),
		GasFeeCap: new(big.Int).SetUint64(gasFeeCap),
		Gas:       gas,
		Nonce:     nonce,
		To:        &tokenContract,
		Data:      callData,
	})

	// AirGap
	if key == nil {
		return tx.Data(), nil
	}

	signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, signedTX)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) TxSendPrepared(ctx *types.RPCContext, tx []byte) (string, error) {
	var signedTx ethTypes.Transaction
	err := signedTx.UnmarshalBinary(tx)

	if err != nil {
		return "", errors.New("cannot unmarshal")
	}

	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, &signedTx)

	return signedTx.Hash().Hex(), err
}

// TX receipt operations

func (e *EVM) TxGetReceipt(ctx *types.RPCContext, tx string) (map[string]interface{}, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx))

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tx_hash":        receipt.TxHash.String(),
		"tx_status":      receipt.Status,
		"tx_type":        receipt.Type,
		"tx_index":       receipt.TransactionIndex,
		"block_number":   receipt.BlockNumber.Uint64(),
		"block_hash":     receipt.BlockHash.String(),
		"gas":            receipt.GasUsed,
		"gas_cumulative": receipt.CumulativeGasUsed,
	}, nil
}

func (e *EVM) GetRPCInfo(ctx *types.RPCContext) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"errors": "",
	}

	height, err := e.GetHeight(ctx)

	if err != nil {
		return nil, err
	}

	block, err := e.GetBlock(ctx, height)

	block.BaseFee()

	if err != nil {
		return nil, err
	}

	gasTipCap, err := e.GetGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	price := 0.0

	gasEstimated := float64(block.GasLimit()*(gasPrice.Uint64()/gwei+gasTipCap.Uint64())/gwei) / float64(gwei)
	calculatedGas := float64(gasMinLimit*(block.BaseFee().Uint64()/gwei+gasTipCap.Uint64())) / float64(gwei) // BSC check
	calculatedGasStr := fmt.Sprintf("%d * (%f + %f) = %f (%f USD)",
		gasMinLimit,
		float64(block.BaseFee().Uint64())/float64(gwei),
		float64(gasTipCap.Uint64())/float64(gwei),
		calculatedGas,
		calculatedGas*price,
	)

	result["name"] = e.Name()
	result["currency"] = e.Currency()
	result["chain_id"], _ = e.GetChainId(ctx)
	result["last_block"] = fmt.Sprintf("%d", height)
	result["gas_limit"] = fmt.Sprintf("%d", block.GasLimit())
	result["gas_tip_cap"] = fmt.Sprintf("%.20f", float64(gasTipCap.Uint64())/float64(gwei))
	result["gas_price"] = fmt.Sprintf("%f gwei (%.20f Ether)", float64(gasPrice.Uint64())/float64(gwei), float64(gasPrice.Uint64())/float64(wei))
	result["gas_base_fee"] = fmt.Sprintf("%.20f", float64(block.BaseFee().Uint64())/float64(gwei))
	result["gas_estimated"] = fmt.Sprintf("%.20f", gasEstimated)
	result["gas_calculated"] = calculatedGasStr
	result["currency_price"] = fmt.Sprintf("%f USD", price)

	/*
		name: Ethereum
		chain_id: 1
		gas_limit: 30000000
		gas_tip_cap: 0.18071816800000001235
		gas_base_fee: 9.66714889600000049086
		gas_estimated: 0.00542154499999999984
		gas_calculated: 21000 * (9.667149 + 0.180718) = 3795.081717
		currency: ETH
		last_block: 16057259
		gas_price: 9.847867 gwei (0.00000000984786706400 Ether)
	*/
	return result, nil
}

func (e *EVM) IsContractAddr(ctx *types.RPCContext, addr string) (bool, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return false, err
	}
	bytecode, err := client.CodeAt(ctx, common.HexToAddress(addr), nil)
	if err != nil {
		log.Fatal(err)
	}

	return len(bytecode) > 0, nil
}

// GetPrice Deprecated
func (e *EVM) GetPrice(ctx *types.RPCContext, contract string) (float64, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}

	isContract, err := e.IsContractAddr(ctx, contract)
	if err != nil {
		return 0, err
	}

	if !isContract {
		return 0, errors.New("address is not contract")
	}

	chainlinkPriceFeedProxy, err := chainlink.NewChainlink(common.HexToAddress(contract), client)
	if err != nil {
		return 0, nil
	}
	decimals, err := chainlinkPriceFeedProxy.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, nil
	}

	roundData, err := chainlinkPriceFeedProxy.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return 0, nil
	}

	return float64(roundData.Answer.Uint64()) / math.Pow(10, float64(decimals)), nil
}

func (e *EVM) ChainLinkGetPrice(ctx *types.RPCContext, contract string) (uint64, uint8, error) {
	client, err := e.GetClient(ctx.NodeId())
	if err != nil {
		return 0, 0, err
	}

	chainlinkPriceFeedProxy, err := chainlink.NewChainlink(common.HexToAddress(contract), client)
	if err != nil {
		return 0, 0, err
	}
	decimals, err := chainlinkPriceFeedProxy.Decimals(&bind.CallOpts{})
	if err != nil {
		return 0, 0, err
	}

	roundData, err := chainlinkPriceFeedProxy.LatestRoundData(&bind.CallOpts{})
	if err != nil {
		return 0, 0, err
	}

	return roundData.Answer.Uint64(), decimals, nil
}

func (e *EVM) GetBlockJson(ctx *types.RPCContext, blockNumber uint64) ([]byte, error) {
	block, err := e.GetBlock(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	return json.Marshal(block)
}

func uint64Pow(base, exp uint64) uint64 {
	result := uint64(1)
	for {
		if exp&1 == 1 {
			result *= base
		}
		exp >>= 1
		if exp == 0 {
			break
		}
		base *= base
	}

	return result
}

func StrToWei(value string) (*big.Int, error) {
	dotIndex := strings.Index(value, `.`)
	weiResult := new(big.Int)
	weiMod := wei

	if dotIndex > 0 {
		floatVal := strings.TrimRight(value[dotIndex+1:], "0")

		if len(floatVal) > 18 {
			return nil, errors.New("max precision is 18")
		}

		floatMod := uint64Pow(10, uint64(len(floatVal)))
		weiMod /= floatMod

		strWei := value[:dotIndex] + floatVal

		_, ok := weiResult.SetString(strWei, 10)

		if !ok {
			return nil, errors.New("cannot parse numeric value")
		}

	} else {
		_, ok := weiResult.SetString(value, 10)

		if !ok {
			return nil, errors.New("cannot parse numeric value")
		}
	}
	weiResult.Mul(weiResult, new(big.Int).SetUint64(weiMod))
	return weiResult, nil
}
