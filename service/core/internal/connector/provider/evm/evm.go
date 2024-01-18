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

package evm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/service/core/internal/connector/client"
	"github.com/censync/soikawallet/service/core/internal/connector/client/evm"
	"github.com/censync/soikawallet/service/core/internal/connector/types/callopts"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/censync/soikawallet/types/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

const (
	ProviderTypeEVM = "evm"
	wei             = uint64(1e18)
	gwei            = uint64(1e9)
	gasMinLimit     = 21000
)

type EVM struct {
	ctx    *types.RPCContext
	client *evm.ClientEVM
}

func NewEVM() *EVM {
	return &EVM{}
}

func (e *EVM) GetType() string {
	return ProviderTypeEVM
}

func (e *EVM) WithClient(ctx *types.RPCContext, evmClient client.Client) (*EVM, error) {
	// panic for sa
	e.client = evmClient.(*evm.ClientEVM)
	e.ctx = ctx
	return e, nil
}

func (e *EVM) GetHeight() (uint64, error) {
	result, err := e.client.Call(e.ctx.Context, "eth_blockNumber", []interface{}{})
	if err != nil {
		return 0, err
	}

	blockNumber, ok := result.(uint64)
	if !ok {
		return 0, errors.New("cannot parse block number")
	}
	return blockNumber, nil
}

func (e *EVM) GetChainId() (*big.Int, error) {
	return nil, nil
}

type rpcTransaction struct {
	tx *ethTypes.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcBlock struct {
	Hash         common.Hash            `json:"hash"`
	Transactions []rpcTransaction       `json:"transactions"`
	UncleHashes  []common.Hash          `json:"uncles"`
	Withdrawals  []*ethTypes.Withdrawal `json:"withdrawals,omitempty"`
}

func (e *EVM) GetBlock(blockNumber uint64) (*ethTypes.Block, error) {
	formattedBlockNum := utils.EncodeUint64(blockNumber)

	response, err := e.client.Call(e.ctx.Context, "eth_getBlockByNumber", []interface{}{formattedBlockNum, true})
	if err != nil {
		return nil, err
	}

	// Decode header and transactions.
	var head *ethTypes.Header
	if err := json.Unmarshal(response.(json.RawMessage), &head); err != nil {
		return nil, err
	}
	// When the block is not found, the API returns JSON null.
	if head == nil {
		return nil, ethereum.NotFound
	}

	var body rpcBlock
	if err := json.Unmarshal(response.(json.RawMessage), &body); err != nil {
		return nil, err
	}

	// Quick-verify transaction and uncle lists. This mostly helps with debugging the server.
	if head.UncleHash == ethTypes.EmptyUncleHash && len(body.UncleHashes) > 0 {
		return nil, errors.New("server returned non-empty uncle list but block header indicates no uncles")
	}
	if head.UncleHash != ethTypes.EmptyUncleHash && len(body.UncleHashes) == 0 {
		return nil, errors.New("server returned empty uncle list but block header indicates uncles")
	}
	if head.TxHash == ethTypes.EmptyTxsHash && len(body.Transactions) > 0 {
		return nil, errors.New("server returned non-empty transaction list but block header indicates no transactions")
	}
	if head.TxHash != ethTypes.EmptyTxsHash && len(body.Transactions) == 0 {
		return nil, errors.New("server returned empty transaction list but block header indicates transactions")
	}
	// Load uncles because they are not included in the block response.
	var uncles []*ethTypes.Header
	if len(body.UncleHashes) > 0 {
		uncles = make([]*ethTypes.Header, len(body.UncleHashes))
		opts := make([]*callopts.CallOpts, len(body.UncleHashes))
		for i := range opts {
			opts[i] = callopts.NewCallOpts(
				"eth_getUncleByBlockHashAndIndex",
				[]interface{}{body.Hash, hexutil.EncodeUint64(uint64(i))},
			)
		}

		batchResp, err := e.client.CallBatch(e.ctx.Context, opts)

		if err != nil {
			return nil, err
		}

		for i := range batchResp.([]interface{}) {
			/*if reqs[i].Error != nil {
				return nil, reqs[i].Error
			}*/
			if uncles[i] == nil {
				return nil, fmt.Errorf("got null header for uncle %d of block %x", i, body.Hash[:])
			}
		}
	}
	// Fill the sender cache of transactions in the block.
	txs := make([]*ethTypes.Transaction, len(body.Transactions))
	/*for i, tx := range body.Transactions {
		if tx.From != nil {
			setSenderFromServer(tx.tx, *tx.From, body.Hash)
		}
		txs[i] = tx.tx
	}*/
	return ethTypes.NewBlockWithHeader(head).WithBody(txs, uncles).WithWithdrawals(body.Withdrawals), nil
}
