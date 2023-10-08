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

package settings

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/service/tui/twidget/tabs"
	"github.com/censync/tview"
)

type pageSettings struct {
	*twidget.BaseFrame
	*state.State

	// vars
	layoutRPCList *tview.Flex
}

func NewPageSettings(state *state.State) *pageSettings {
	layout := tview.NewFlex()

	return &pageSettings{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageSettings) FuncOnShow() {
	tabs := tabs.NewTabs().
		AddItem("Application", p.tabApp()).
		AddItem("Labels", p.tabLabels()).
		AddItem("RPC", p.tabNodes())
	p.BaseLayout().AddItem(tabs, 0, 1, false)

}
