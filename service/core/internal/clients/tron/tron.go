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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/censync/soikawallet/service/core/internal/apiclient"
	"github.com/censync/soikawallet/service/core/internal/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"net"
	"net/http"
	"time"
)

const unit = 1000000

type Tron struct {
	*types.BaseNetwork
}

func NewTron(baseNetwork *types.BaseNetwork) *Tron {
	return &Tron{BaseNetwork: baseNetwork}
}

/*func (n *Tron) AddressHex(pub *ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(*pub).Hex()
	address = "41" + address[2:]
	return address
}*/

func (t *Tron) getClient(nodeId uint32) *apiclient.ApiClient {
	return &apiclient.ApiClient{
		Host: t.DefaultRPC().Endpoint(),
		Client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).DialContext,
				DisableCompression:    false,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     true,
				MaxIdleConnsPerHost:   -1,
			},
		},
	}
}

func (t *Tron) Address(pub *ecdsa.PublicKey) string {
	addr := crypto.PubkeyToAddress(*pub).Hex()
	addr = "41" + addr[2:]
	addb, _ := hex.DecodeString(addr)
	hash1 := s256(s256(addb))
	secret := hash1[:4]
	for _, v := range secret {
		addb = append(addb, v)
	}
	return base58.Encode(addb)
}

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func (t *Tron) getHeight(ctx *types.RPCContext) (uint64, error) {
	var res struct {
		Header struct {
			RawData struct {
				Number uint64 `json:"number"`
			} `json:"raw_data"`
		} `json:"block_header"`
	}
	method := `/wallet/getnowblock`

	err := t.getClient(ctx.NodeId()).Do("POST", method, nil, &res)
	if err != nil {
		return 0, err
	}

	return res.Header.RawData.Number, nil
}

func (t *Tron) GetBalance(ctx *types.RPCContext) (float64, error) {
	var res struct {
		Balance uint64 `json:"balance"`
	}
	method := `/wallet/getaccount`
	req := struct {
		Address string `json:"address"`
		Visible bool   `json:"visible"`
	}{
		Address: ctx.CurrentAccount(),
		Visible: true,
	}
	err := t.getClient(ctx.NodeId()).Do("POST", method, &req, &res)
	if err != nil {
		return 0, err
	}

	return float64(res.Balance) / float64(unit), nil
}

func (t *Tron) GetTokenBalance(ctx *types.RPCContext, contract string, decimals int) (*big.Float, error) {
	return new(big.Float), nil
}

// gettransactioninfobyid
func (t *Tron) TxGetReceipt(ctx *types.RPCContext, tx string) (map[string]interface{}, error) {
	var res struct {
		Ret struct {
			ContractRet string `json:"contractRet"`
		} `json:"ret"`
		TxID    string `json:"txID"`
		RawData struct {
			Contract []struct {
				Type      string `json:"type"`
				Parameter struct {
					Value struct {
						OwnerAddress string  `json:"owner_address"`
						ToAddress    string  `json:"to_address"`
						Amount       float64 `json:"amount"`
					} `json:"value"`
				} `json:"parameter"`
			} `json:"contract"`
		} `json:"raw_data"`
	}
	method := `/wallet/gettransactionbyid`
	req := struct {
		Value string `json:"value"`
	}{
		Value: tx,
	}
	err := t.getClient(ctx.NodeId()).Do("POST", method, &req, &res)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tx_hash":        res.TxID,
		"tx_status":      res.Ret.ContractRet,
		"tx_type":        res.RawData.Contract[0].Type,
		"tx_index":       -1,
		"block_number":   -1,
		"block_hash":     -1,
		"gas":            -1,
		"gas_cumulative": -1,
	}, nil
}

func (t *Tron) GetRPCInfo(ctx *types.RPCContext) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	result["name"] = t.Name()
	result["currency"] = t.Currency()
	result["last_block"], _ = t.getHeight(ctx)

	return result, nil
}
