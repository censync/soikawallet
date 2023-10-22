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

package dto

import (
	mhda "github.com/censync/go-mhda"
	resp "github.com/censync/soikawallet/api/responses"
)

type ConnectDTO struct {
	InstanceId string
	Origin     string
	RemoteAddr string
}

type ExecuteRPCRequestDTO struct {
	InstanceId string
	Origin     string
	RemoteAddr string

	ChainKey mhda.ChainKey
	Method   string
	Params   []interface{}
}

type ResponseAcceptDTO struct {
	InstanceId string
	Chains     []*resp.ChainInfo
}

type ResponseProxyCallDTO struct {
	InstanceId string
	Data       []byte
}

type ResponseRejectDTO struct {
	InstanceId string
	RemoteAddr string
}

type RequestAccountsDTO struct {
	InstanceId string
	Origin     string
	ChainKey   mhda.ChainKey
}

type RequestCallGetBlockByNumberDTO struct {
	InstanceId string
	Origin     string
	ChainKey   mhda.ChainKey
	Method     string
	Params     []interface{}
}
