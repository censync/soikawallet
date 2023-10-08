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

package state

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/core"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/twidget/extpages"
)

const (
	ModeIdle = iota
	ModeWithAirGap
	ModeWithoutAirGap

	StatusIdle = iota
	StateInitAirGap
	StateInitLocal
	StateAwaitQR
)

type State struct {
	uiEvents *events.EventBus
	w3Events *events.EventBus

	walletMode uint8
	status     uint8

	isInitialised bool
	tr            *i18n.Translator
	pages         *extpages.ExtPages
}

func InitState(uiEvents, w3Events *events.EventBus, tr *i18n.Translator) *State {
	return &State{
		uiEvents:   uiEvents,
		w3Events:   w3Events,
		walletMode: ModeIdle,
		status:     StatusIdle,
		tr:         tr,
	}
}

func (s *State) Emit(event events.EventType, data interface{}) {
	s.uiEvents.Emit(event, data)
}

func (s *State) EmitW3(event events.EventType, data interface{}) {
	s.w3Events.Emit(event, data)
}

func (s *State) WalletMode() uint8 {
	return s.walletMode
}

func (s *State) Status() uint8 {
	return s.status
}

func (s *State) API() core.CoreAdapter {
	return core.API()
}

func (s *State) Tr() *i18n.Translator {
	return s.tr
}

func (s *State) SetStatus(status uint8) {
	s.status = status
}

func (s *State) SetWalletMode(mode uint8) {
	s.walletMode = mode
}

// Pages
func (s *State) SetPages(pages *extpages.ExtPages) {
	s.pages = pages
}

func (s *State) Pages() *extpages.ExtPages {
	return s.pages
}

func (s *State) SwitchToPage(page string, args ...interface{}) {

	// TODO: Change to channel based uiEvents
	s.pages.SwitchToPage(page, args...)
	//s.Emit(events.EventDrawForce, nil)
}

// Current pages

func (s *State) Params() []interface{} {
	return s.pages.Current().Params()
}
