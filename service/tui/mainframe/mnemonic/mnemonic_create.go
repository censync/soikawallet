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
	"github.com/censync/soikawallet/service/tui/util/clipboard"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

type pageInitMnemonic struct {
	*twoframes.BaseFrame
	*state.State

	// ui
	inputMnemonic *tview.Form

	// read-only
	entropyList []string

	// vars
	selectedMnemonicEntropy  int
	selectedMnemonicLanguage string
	mnemonic                 string
}

func NewPageInitMnemonic(state *state.State) *pageInitMnemonic {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageInitMnemonic{
		State:       state,
		entropyList: []string{"128", "160", "192", "224", "256"},
		BaseFrame:   twoframes.NewBaseFrame(layout),
	}
}

func (p *pageInitMnemonic) FuncOnShow() {
	p.inputMnemonic = tview.NewForm().
		SetHorizontal(true)

	p.inputMnemonic.SetItemPadding(4).
		SetBorder(true).
		SetTitle(` ` + p.Tr().T("ui.label", "mnemonic") + ` `).
		SetTitleAlign(tview.AlignLeft)

	inputPassword := tview.NewInputField().
		SetMaskCharacter('*')
	inputPassword.SetTitle(` ` + p.Tr().T("ui.label", "passphrase") + ` `).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	inputDropDownEntropy := tview.NewDropDown().
		SetLabel(p.Tr().T("ui.label", "entropy")).
		SetFieldWidth(5).
		SetOptions(p.entropyList, func(option string, optionIndex int) {
			p.selectedMnemonicEntropy, _ = strconv.Atoi(option)
			p.actionMnemonicUpdate()
		}).
		SetCurrentOption(len(p.entropyList) - 1)

	inputDropDownLanguage := tview.NewDropDown().
		SetLabel(p.Tr().T("ui.label", "language")).
		SetFieldWidth(10).
		SetOptions([]string{"english"}, func(option string, optionIndex int) {
			if p.selectedMnemonicLanguage != option {
				p.selectedMnemonicLanguage = option
				p.actionMnemonicUpdate()
			}
		}).
		SetCurrentOption(0)

	formMnemonicConfig := tview.NewForm().
		SetHorizontal(true).
		AddFormItem(inputDropDownEntropy).
		AddFormItem(inputDropDownLanguage)

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(func() {
			instanceId, err := core.API().Init(&dto.InitWalletDTO{
				Mnemonic:   p.mnemonic,
				Passphrase: inputPassword.GetText(),
			})
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot init wallet: %s", err))
			} else {
				p.Emit(events.EventUpdateCurrencies, nil)
				clipboard.Clear()
				p.Emit(events.EventWalletInitialized, instanceId)
				p.SwitchToPage(pages.CreateWallets)
			}
		})
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

	layoutOptions := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "generate_mnemonic"), func() {
			p.actionMnemonicUpdate()
		}).
		AddButton(p.Tr().T("ui.button", "copy_to_clipboard"), func() {
			err := clipboard.CopyToClipboard(p.mnemonic)
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot copy to clipboard: %s", err))
			}
		})

	layoutMnemonicForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.inputMnemonic, 0, 6, false).
		AddItem(inputPassword, 3, 1, false).
		AddItem(formMnemonicConfig, 3, 1, false).
		AddItem(layoutOptions, 3, 1, false)

	p.BaseLayout().AddItem(layoutMnemonicForm, 0, 3, false).
		AddItem(layoutWizard, 35, 1, false)

}

func (p *pageInitMnemonic) actionMnemonicUpdate() {
	var err error
	//p.inputMnemonic.SetText(``, false)
	p.inputMnemonic.Clear(false)

	p.mnemonic, err = core.API().GenerateMnemonic(&dto.GenerateMnemonicDTO{
		BitSize:  p.selectedMnemonicEntropy,
		Language: p.selectedMnemonicLanguage,
	})

	if err != nil {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot generate mnemonic: %s", err))
		return
	}

	//p.inputMnemonic.SetText(mnemonic, false)

	mnemonicList := strings.Split(p.mnemonic, ` `)

	for index := range mnemonicList {
		mnemonicInput := tview.NewInputField().
			SetFieldWidth(15).
			SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
				return false
			}).
			SetLabel(fmt.Sprintf("%2d.", index+1)).
			SetText(mnemonicList[index])

		p.inputMnemonic.AddFormItem(mnemonicInput)
	}
}

func (p *pageInitMnemonic) FuncOnHide() {
	p.mnemonic = ``
	p.BaseLayout().Clear()
}
