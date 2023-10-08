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
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageInitWallet struct {
	*twidget.BaseFrame
	*state.State
}

func NewPageInitWallet(state *state.State) *pageInitWallet {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageInitWallet{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageInitWallet) FuncOnShow() {
	btnWalletCreate := tview.NewButton(p.Tr().T("ui.button", "wallet_create"))

	btnWalletCreate.SetSelectedFunc(func() {
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(pages.MnemonicInit)
	})
	btnWalletRestore := tview.NewButton(p.Tr().T("ui.button", "wallet_restore"))

	btnWalletRestore.SetSelectedFunc(func() {
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(pages.MnemonicRestore)
	})

	layoutButtons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(btnWalletCreate, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(btnWalletRestore, 0, 1, false).
		AddItem(nil, 0, 1, false)

	layoutButtons.SetBorderPadding(0, 0, 10, 10)

	p.BaseLayout().AddItem(layoutButtons, 3, 1, false)
}
