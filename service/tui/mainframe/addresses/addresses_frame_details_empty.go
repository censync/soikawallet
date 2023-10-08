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

package addresses

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
)

type frameAddressesDetailsEmpty struct {
	layout *tview.Flex
	*state.State
}

func newFrameAddressesDetailsEmpty(state *state.State) *frameAddressesDetailsEmpty {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameAddressesDetailsEmpty{
		State:  state,
		layout: layout,
	}
}

func (f *frameAddressesDetailsEmpty) Layout() *tview.Flex {
	label := tview.NewTextView().
		SetText("Account or address not selected")
	label.SetBorderPadding(0, 0, 8, 8)
	f.layout.AddItem(nil, 0, 1, false).
		AddItem(label, 0, 1, false).
		AddItem(nil, 0, 1, false)
	return f.layout
}
