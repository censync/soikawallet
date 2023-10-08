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
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package tui

import (
	"github.com/censync/soikawallet/config"
	"github.com/censync/soikawallet/service/tui/api_web3"
	"github.com/censync/soikawallet/service/tui/events"
	"sync"
)

var provider *Service

type IServiceProvider interface {
	Start() error
	Stop()
}

type Service struct {
	uiEvents              *events.EventBus
	w3Events              *events.EventBus
	web3ConnectionService *api_web3.Web3Connection
	tuiService            *Tui
}

func NewTUIServiceProvider(cfg *config.Config, wg *sync.WaitGroup) *Service {
	uiEvents := events.NewEventBus()
	w3Events := events.NewEventBus()
	return &Service{
		web3ConnectionService: api_web3.NewWeb3Connection(cfg, wg, uiEvents, w3Events),
		tuiService:            NewTui(cfg, wg, uiEvents, w3Events),
	}
}

func (p *Service) Web3Connection() IServiceProvider {
	return p.web3ConnectionService
}

func (p *Service) UI() IServiceProvider {
	return p.tuiService
}
