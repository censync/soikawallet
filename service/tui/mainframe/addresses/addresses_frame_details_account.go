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

package addresses

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
)

const (
	AccountLabel = 1
	AddressLabel = 2
)

type frameAddressesDetailsAccount struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedAccount *resp.AccountResponse

	selectedLabelIndex uint32
}

func newAddressesFrameDetailsAccount(state *state.State, accountPath *resp.AccountResponse) *frameAddressesDetailsAccount {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameAddressesDetailsAccount{
		State:           state,
		layout:          layout,
		selectedAccount: accountPath,
	}
}

func (f *frameAddressesDetailsAccount) Layout() *tview.Flex {
	label := tview.NewTextView().
		SetText("Account selected")

	label.SetBorderPadding(0, 0, 8, 8)

	accountLabels := f.API().GetAccountLabels()

	inputSelectLabel := tview.NewDropDown().
		SetLabel("Select label")

	for index, title := range accountLabels {
		inputSelectLabel.AddOption(title, func() {
			f.selectedLabelIndex = index
		})
	}

	inputSelectLabel.AddOption(" [ add label ] ", func() {
		f.SwitchToPage(pages.Settings)
	})

	formAccountDesc := tview.NewForm().
		AddFormItem(inputSelectLabel).
		AddButton("set label", f.actionSetLabel).
		AddButton("Remove label", f.actionRemoveLabel)

	f.layout.AddItem(label, 1, 1, false).
		AddItem(formAccountDesc, 0, 1, false)
	return f.layout
}

func (f *frameAddressesDetailsAccount) actionSetLabel() {
	err := f.API().SetLabelLink(&dto.SetLabelLinkDTO{
		LabelType: AccountLabel,
		Index:     f.selectedLabelIndex,
		// TODO: Finish
		//Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(events.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot set label: %s", err))
	}
}

func (f *frameAddressesDetailsAccount) actionRemoveLabel() {
	err := f.API().RemoveLabelLink(&dto.RemoveLabelLinkDTO{
		LabelType: AccountLabel,
		// TODO: Finish
		// Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(events.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot remvove label: %s", err))
	}
}
