package evm

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/service/wallet/internal/oracle/chainlink"
	"github.com/censync/soikawallet/service/wallet/internal/token/erc20"
	"github.com/censync/soikawallet/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/big"
)

const (
	wei         = uint64(1e18)
	gwei        = uint64(1e9)
	gasMinLimit = 21000
)

type EVM struct {
	*types.BaseNetwork
	client map[uint32]*ethclient.Client
}

func NewEVM(baseNetwork *types.BaseNetwork) *EVM {
	return &EVM{BaseNetwork: baseNetwork, client: map[uint32]*ethclient.Client{}}
}

func (e *EVM) getClient(nodeId uint32) (*ethclient.Client, error) {
	var err error
	if e.client[nodeId] != nil {
		return e.client[nodeId], nil
	} else {
		e.client[nodeId], err = ethclient.Dial(e.DefaultRPC().Endpoint())
		return e.client[nodeId], err
	}
}

func (e *EVM) Address(pub *ecdsa.PublicKey) string {
	return crypto.PubkeyToAddress(*pub).Hex()
}

func (e *EVM) getHeight(ctx *types.RPCContext) (uint64, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}
	return client.BlockNumber(ctx)
}

func (e *EVM) getBlock(ctx *types.RPCContext, blockNumber uint64) (*ethTypes.Block, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}
	return client.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
}

func (e *EVM) GetBalance(ctx *types.RPCContext) (float64, error) {
	client, err := e.getClient(ctx.NodeId())
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
	client, err := e.getClient(ctx.NodeId())
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
	client, err := e.getClient(ctx.NodeId())
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
	client, err := e.getClient(ctx.NodeId())
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

func (e *EVM) ApproveToken(ctx *types.RPCContext, contract string, value float64) (string, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}
	instance, err := erc20.NewErc20(common.HexToAddress(contract), client)
	if err != nil {
		return ``, err
	}
	weiValue := big.NewInt(int64(value * float64(wei)))
	signedTX, err := instance.Approve(&bind.TransactOpts{}, common.HexToAddress(ctx.CurrentAccount()), weiValue)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) getGasPrice(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.SuggestGasPrice(ctx)
}

func (e *EVM) getNonce(ctx *types.RPCContext, account string) (uint64, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}

	return client.PendingNonceAt(ctx, common.HexToAddress(account))
}

func (e *EVM) getChainId(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.ChainID(ctx)
}

func (e *EVM) getGasTipCap(ctx *types.RPCContext) (*big.Int, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return nil, err
	}

	return client.SuggestGasTipCap(ctx)
}

func (e *EVM) getGasEstimate(ctx *types.RPCContext, msg *ethereum.CallMsg) (uint64, error) {
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return 0, err
	}

	return client.EstimateGas(ctx, *msg)
}

func (e *EVM) GetGasBaseTx(ctx *types.RPCContext) (map[string]float64, error) {
	result := map[string]float64{
		"base_fee":     0,
		"priority_fee": 0,
		"units":        21000,
		"limit":        0,
		"currency":     0,
	}

	height, err := e.getHeight(ctx)

	if err != nil {
		return result, err
	}

	block, err := e.getBlock(ctx, height)

	if err != nil {
		return result, err
	}

	gasLimit := block.GasLimit()
	result["limit"] = float64(gasLimit)

	baseFee := block.BaseFee()
	if baseFee != nil {
		result["base_fee"] = float64(baseFee.Int64()) / float64(gwei)
	}

	gasTipCap, err := e.getGasTipCap(ctx)
	if err != nil {
		return result, err
	}
	if gasTipCap != nil {
		result["priority_fee"] = float64(gasTipCap.Int64()) / float64(gwei)
	}
	return result, nil
}

func (e *EVM) TxSendBase(ctx *types.RPCContext, to string, value float64, key *ecdsa.PrivateKey) (txId string, err error) {
	chainId, err := e.getChainId(ctx)

	if err != nil {
		return "", err
	}

	height, err := e.getHeight(ctx)

	if err != nil {
		return "", err
	}

	block, err := e.getBlock(ctx, height)

	if err != nil {
		return "", err
	}

	gasTipCap, err := e.getGasTipCap(ctx)

	if err != nil {
		return "", err
	}

	nonce, err := e.getNonce(ctx, ctx.CurrentAccount())

	if err != nil {
		return "", err
	}

	gasLimit := block.GasLimit()

	addrTo := common.HexToAddress(to)
	weiValue := big.NewInt(int64(value * float64(wei)))
	maxGas := big.NewInt(int64(gasLimit - 1000000))

	// both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified

	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasFeeCap: maxGas,
		GasTipCap: gasTipCap,
		Gas:       gasTipCap.Uint64() + 10000000, // 10000000
		Nonce:     nonce,
		To:        &addrTo,
		Value:     weiValue,
		Data:      nil,
	})

	signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

	if err != nil {
		return ``, nil
	}

	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, signedTX)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) TxSendToken(ctx *types.RPCContext, to string, value float64, token *types.TokenConfig, key *ecdsa.PrivateKey) (txId string, err error) {
	chainId, err := e.getChainId(ctx)

	if err != nil {
		return "", err
	}

	/*height, err := e.getHeight(ctx)

	if err != nil {
		return "", err
	}

	block, err := e.getBlock(ctx, height)

	gasLimit := block.GasLimit()*/

	if err != nil {
		return "", err
	}

	gasTipCap, err := e.getGasTipCap(ctx)

	if err != nil {
		return "", err
	}

	nonce, err := e.getNonce(ctx, ctx.CurrentAccount())

	if err != nil {
		return "", err
	}

	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return ``, err
	}
	allowance, err := e.GetTokenAllowance(ctx, token.Contract(), to)

	if err != nil {
		return "", err
	}

	if allowance == 0 {
		return ``, errors.New("not approved")
	}

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	addrFrom := common.HexToAddress(to)
	paddedAddress := common.LeftPadBytes(addrFrom.Bytes(), 32)

	addrTo := common.HexToAddress(to)

	weiValue := big.NewInt(int64(value * float64(wei)))

	paddedAmount := common.LeftPadBytes(weiValue.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	tokenContract := common.HexToAddress(token.Contract())

	dataTx := ethereum.CallMsg{
		To:   &tokenContract,
		Data: data,
	}

	gasEstimate, err := e.getGasEstimate(ctx, &dataTx)

	if err != nil {
		return ``, err
	}

	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasFeeCap: big.NewInt(int64(gasEstimate)),
		GasTipCap: gasTipCap,
		Gas:       gasTipCap.Uint64() + 10000000, // 10000000
		Nonce:     nonce,
		To:        &addrTo,
		Value:     weiValue,
		Data:      data,
	})

	signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

	if err != nil {
		return ``, err
	}

	err = client.SendTransaction(ctx, signedTX)

	return signedTX.Hash().Hex(), err
}

func (e *EVM) TxPrepare(ctx *types.RPCContext, to string, value float64) (interface{}, error) {
	/*chainId, err := e.getChainId(ctx)

	if err != nil {
		return nil, err
	}

	height, err := e.getHeight(ctx)

	if err != nil {
		return nil, err
	}

	block, err := e.getBlock(ctx, height)

	if err != nil {
		return nil, err
	}

	gasTipCap, err := e.getGasTipCap(ctx)

	if err != nil {
		return nil, err
	}

	nonce, err := e.getNonce(ctx, to)

	if err != nil {
		return nil, err
	}

	gasLimit := block.GasLimit()

	addrTo := common.HexToAddress(to)
	weiValue := big.NewInt(int64(float64(wei) * value))
	maxGas := big.NewInt(int64(gasLimit - 21000))
	tx := ethTypes.NewTx(&ethTypes.DynamicFeeTx{
		ChainID:   chainId,
		GasFeeCap: maxGas,
		GasTipCap: gasTipCap,
		//Gas:       maxGas * gasTipCap.Uint64(), // 10000000
		Nonce: nonce,
		To:    &addrTo,
		Value: weiValue,
		Data:  nil,
	})
	return tx.Data(), nil*/
	/*
		signedTX, err := ethTypes.SignTx(tx, ethTypes.LatestSignerForChainID(chainId), key)

		if err != nil {
			return ``, nil
		}

		client, err := e.getClient(ctx.NodeId())
		if err != nil {
			return ``, err
		}

		err = client.SendTransaction(context.Background(), signedTX)

		return signedTX.Hash().Hex(), err*/
	return nil, nil
}

func (e *EVM) TxSend(ctx *context.Context) error {
	return nil
}

func (e *EVM) TxGetReceipt(ctx *types.RPCContext, tx string) (map[string]interface{}, error) {
	client, err := e.getClient(ctx.NodeId())
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

	height, err := e.getHeight(ctx)

	if err != nil {
		return nil, err
	}

	block, err := e.getBlock(ctx, height)

	block.BaseFee()

	if err != nil {
		return nil, err
	}

	gasTipCap, err := e.getGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.getGasPrice(ctx)
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
	result["chain_id"], _ = e.getChainId(ctx)
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
	client, err := e.getClient(ctx.NodeId())
	if err != nil {
		return false, err
	}
	bytecode, err := client.CodeAt(ctx, common.HexToAddress(addr), nil)
	if err != nil {
		log.Fatal(err)
	}

	return len(bytecode) > 0, nil
}

func (e *EVM) GetPrice(ctx *types.RPCContext, contract string) (float64, error) {
	client, err := e.getClient(ctx.NodeId())
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
