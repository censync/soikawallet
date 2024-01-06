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

package init_wallet

import (
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageInitMode struct {
	*twoframes.BaseFrame
	*state.State
}

func NewPageInitMode(state *state.State) *pageInitMode {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	layout.SetBorderPadding(5, 5, 5, 5)

	return &pageInitMode{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageInitMode) FuncOnShow() {
	btnModeAirGap := tview.NewButton("Use AirGap [[green]Recommended[white]]")
	btnModeAirGap.SetSelectedFunc(func() {
		p.SetWalletMode(state.ModeWithAirGap)
		p.SetStatus(state.StateInitAirGap)
		p.SwitchToPage(pages.AirGapShow)
	})

	btnModeLocal := tview.NewButton("Do not use AirGap [[red]less secure[white]]")
	btnModeLocal.SetSelectedFunc(func() {
		p.SetWalletMode(state.ModeWithoutAirGap)
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(pages.SelectInitWallet)
	})

	labelAirGap := tview.NewTextView().SetText(p.Tr().T("ui.label", "splash_option_airgap"))
	labelLocal := tview.NewTextView().SetText(p.Tr().T("ui.label", "splash_option_no_airgap"))

	labelAirGap.SetBorderPadding(0, 1, 2, 2)
	labelLocal.SetBorderPadding(0, 1, 2, 2)

	layoutTop := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(labelAirGap, 0, 1, false).
		AddItem(labelLocal, 0, 1, false)

	layoutButtons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(btnModeAirGap, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(btnModeLocal, 0, 1, false).
		AddItem(nil, 0, 1, false)

	layoutButtons.SetBorderPadding(0, 0, 10, 10)

	p.BaseLayout().AddItem(layoutTop, 2, 1, false).
		AddItem(layoutButtons, 3, 1, false)
}
