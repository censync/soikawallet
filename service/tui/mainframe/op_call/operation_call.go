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

package op_call

import (
	"fmt"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/formtextview"
	"github.com/censync/twidget/twoframes"
)

type pageOperationCall struct {
	*twoframes.BaseFrame
	*state.State

	paramSelectedAddr *responses.AddressResponse

	// ui
	// input
	inputAddrSender   *formtextview.FormTextView
	inputAddrContract *tview.InputField
	inputAddrMethod   *tview.InputField

	// wizard
	layoutArgsEntriesForm *tview.Flex
	argCount              int
}

func NewPageOperationCall(state *state.State) *pageOperationCall {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageOperationCall{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
		argCount:  3,
	}
}

func (p *pageOperationCall) Layout() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Settings
	p.inputAddrSender = formtextview.NewFormTextView("")
	p.inputAddrSender.SetLabel("Sender ")

	p.inputAddrContract = tview.NewInputField().
		SetFieldWidth(44)
	p.inputAddrContract.SetLabel("Contract")

	p.inputAddrMethod = tview.NewInputField().
		SetFieldWidth(32)
	p.inputAddrMethod.SetLabel("Method")

	layoutInputForm := tview.NewForm().
		SetHorizontal(false)

	layoutInputForm.AddFormItem(p.inputAddrSender).
		AddFormItem(p.inputAddrContract).
		AddFormItem(p.inputAddrMethod)

	// Options
	layoutOptionsForm := tview.NewForm().
		SetHorizontal(false)

	layoutSettings := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layoutSettings.
		AddItem(layoutInputForm, 60, 1, false).
		AddItem(layoutOptionsForm, 0, 1, false)

	layoutOptionsForm.
		AddButton("Load ABI", func() {
		}).
		AddButton("Load prepared", func() {
		}).
		AddButton("Save prepared", func() {
		}).
		AddButton("add", func() {
			p.argCount++
		}).
		AddButton("remove", func() {
			if p.argCount > 1 {
				p.argCount--
			}
		})

	p.layoutArgsEntriesForm = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutArgsEntriesForm.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(`Builder`)

	layout.
		AddItem(layoutSettings, 7, 1, false).
		AddItem(p.layoutArgsEntriesForm, 0, 1, false)
	return layout
}

func (p *pageOperationCall) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Sender address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.paramSelectedAddr = p.Params()[0].(*responses.AddressResponse)

	p.actionUpdateForm()
}

func (p *pageOperationCall) actionUpdateForm() {
	p.layoutArgsEntriesForm.Clear()

	p.inputAddrSender.SetText(p.paramSelectedAddr.Address)

	for argIndex := 0; argIndex < p.argCount; argIndex++ {
		labelArgsForm := tview.NewForm().
			SetHorizontal(true).
			AddDropDown(
				"type",
				[]string{" Select ▼ ",
					"int", "int8", "int16", "int32", "int64", "int128", "int256",
					"uint", "uint8", "uint16", "uint32", "uint64", "uint128", "uint256",
					"address", "string", "bool", "byte"},
				0,
				func(option string, optionIndex int) {

				}).
			AddInputField("value", "", 20, nil, nil).
			AddButton("Add row", func() {

			}).
			AddButton("Remove row", func() {

			})

		p.layoutArgsEntriesForm.AddItem(labelArgsForm, 3, 1, false)
	}
}

func (p *pageOperationCall) actionOperationCall() {

}

func (p *pageOperationCall) FuncOnHide() {
	p.paramSelectedAddr = nil
	p.BaseLayout().Clear()
}
