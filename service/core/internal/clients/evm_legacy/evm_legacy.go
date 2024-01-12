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

package evm_legacy

import (
	"crypto/ecdsa"
	"github.com/censync/soikawallet/service/core/internal/clients/evm_base"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type EVMLegacy struct {
	*evm_base.EVM
}

func NewEVMLegacy(baseNetwork *types.BaseNetwork) *EVMLegacy {
	return &EVMLegacy{EVM: evm_base.NewEVM(baseNetwork)}
}

func (e *EVMLegacy) TxSendBase(ctx *types.RPCContext, to string, value string, gas, gasFeeCap, _ uint64, key *ecdsa.PrivateKey) (interface{}, error) {
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
	weiValue, err := evm_base.StrToWei(value)

	if err != nil {
		return 0, err
	}

	txData = &ethTypes.LegacyTx{
		GasPrice: new(big.Int).SetUint64(gasFeeCap), // base price
		Gas:      gas,
		Nonce:    nonce,
		To:       &addrTo,
		Value:    weiValue,
		Data:     nil,
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

func (e *EVMLegacy) TxSendToken(ctx *types.RPCContext, to, value string, token *types.TokenConfig, gas, _, gasFeeCap uint64, key *ecdsa.PrivateKey) (interface{}, error) {
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
	weiAmount, err := evm_base.StrToWei(value)

	if err != nil {
		return 0, err
	}

	callData := evm_base.GasCalcPrepared("transfer", addrTo, weiAmount)

	tokenContract := common.HexToAddress(token.Contract())

	txData = &ethTypes.LegacyTx{
		GasPrice: new(big.Int).SetUint64(gasFeeCap), // base price
		Gas:      gas,
		Nonce:    nonce,
		To:       &tokenContract,
		Data:     callData,
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
