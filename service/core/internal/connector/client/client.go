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

package client

import (
	"context"
	"github.com/censync/soikawallet/service/core/internal/connector/types/callopts"
)

type Client interface {
	Index() uint32
	Dial() error
	StartSubscription(string, []interface{}) (<-chan interface{}, <-chan struct{}, error)
	Call(context.Context, string, []interface{}) (interface{}, error)
	CallOpts(context.Context, *callopts.CallOpts) (interface{}, error)
	CallBatch(context.Context, []*callopts.CallOpts) (interface{}, error)
	IsReady() bool
	Stop()
	Disconnect()
}
