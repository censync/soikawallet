// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package api_web3

import "encoding/json"

// Request codes
const (
	reqCodePing = iota + 100
	reqCodeConnect
	reqCodeRequestAccounts
	reqCodeGetBalance
	reqCodeProxyCall
)

// Response codes
const (
	respCodePong uint16 = iota + 100
	respCodeConnectionAccepted
	respCodeConnectionRejected
	respCodeGetAccounts
	respCodeProxyCall
	respCodeError      = 400
	respCodeErrorFatal = 501
)

type RPCMessageHeader struct {
	Version uint8  `json:"_v"`
	Id      string `json:"_id"`
	Type    uint16 `json:"type"`
}

type RPCMessageReq struct {
	RPCMessageHeader
	Data interface{} `json:"data"`
}

type RPCMessageResp struct {
	*RPCMessageHeader
	Data interface{} `json:"data"`
}

func (r *RPCMessageResp) toJSON() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return []byte("{}")
	}
	return data
}
