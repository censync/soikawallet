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

package mnemonic

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/core"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
	"os"
	"strings"
)

type pageRestoreMnemonic struct {
	*twidget.BaseFrame
	*state.State

	// ui
	inputMnemonic *tview.TextArea
	inputPassword *tview.InputField
}

func NewPageRestoreMnemonic(state *state.State) *pageRestoreMnemonic {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageRestoreMnemonic{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageRestoreMnemonic) FuncOnShow() {
	p.inputMnemonic = tview.NewTextArea()
	p.inputMnemonic.SetTitle(p.Tr().T("ui.label", "mnemonic")).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	p.inputPassword = tview.NewInputField().
		SetMaskCharacter('*')
	p.inputPassword.SetTitle(p.Tr().T("ui.label", "passphrase")).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	// env
	envMnemonic, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC")
	if ok {
		p.inputMnemonic.SetText(strings.TrimSpace(envMnemonic), true)
		_ = os.Unsetenv("SOIKAWALLET_MNEMONIC")
	}

	envMnemonicPassphrase, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC_PASSPHRASE")
	if ok {
		p.inputPassword.SetText(strings.TrimSpace(envMnemonicPassphrase))
		_ = os.Unsetenv("SOIKAWALLET_MNEMONIC_PASSPHRASE")
	}

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(p.actionRestoreWithMnemonic)

	btnBack := tview.NewButton(p.Tr().T("ui.button", "back")).
		SetSelectedFunc(func() {
			p.SwitchToPage(p.Pages().GetPrevious())
		})

	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnNext, 3, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(btnBack, 3, 1, false)

	layoutWizard.SetBorderPadding(1, 1, 2, 2)

	layoutRestoreForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.inputMnemonic, 10, 1, false).
		AddItem(p.inputPassword, 3, 1, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutRestoreForm, 0, 3, false).
		AddItem(layoutWizard, 35, 1, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageRestoreMnemonic) actionRestoreWithMnemonic() {
	instanceId, err := core.API().Init(&dto.InitWalletDTO{
		Mnemonic:   p.inputMnemonic.GetText(),
		Passphrase: p.inputPassword.GetText(),
	})
	if err != nil {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot restore wallet: %s", err))
	} else {
		//p.SetWallet(core)
		p.Emit(events.EventUpdateCurrencies, nil)
		clipboard.Clear()
		p.Emit(events.EventWalletInitialized, instanceId)
		p.SwitchToPage(pages.CreateWallets)
	}
}
