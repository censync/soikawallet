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

//go:build testnet

package about

import (
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
	"os"
	"strings"
)

type pageAgreement struct {
	*twoframes.BaseFrame
	*state.State

	isAgreementAccepted bool
}

func NewPageAgreement(state *state.State) *pageAgreement {
	var isAgreementAccepted bool

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(2, 2, 3, 3)

	envAgreementAccepted, ok := os.LookupEnv("SOIKAWALLET_AGREEMENT_ACCEPTED")

	if ok && strings.ToLower(envAgreementAccepted) == "true" {
		isAgreementAccepted = true
	}

	return &pageAgreement{
		BaseFrame:           twoframes.NewBaseFrame(layout),
		State:               state,
		isAgreementAccepted: isAgreementAccepted,
	}
}

func (p *pageAgreement) FuncOnShow() {
	if p.isAgreementAccepted {
		p.SwitchToPage(pages.SelectInitWallet)
	}

	viewTermsOfUse := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorRed).
		SetText(p.Tr().T("ui.label", "terms_of_use_testnet"))

	viewTermsOfUse.SetBorder(true).
		SetTitle(" Terms of use ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	viewPrivacyPolicy := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorRed).
		SetText(p.Tr().T("ui.label", "privacy_policy_testnet"))

	viewPrivacyPolicy.SetBorder(true).
		SetTitle(" Privacy policy ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	btnAccept := tview.NewButton("Accept").
		SetSelectedFunc(func() {
			p.SwitchToPage(pages.SelectInitWallet)
		})

	btnDecline := tview.NewButton("Decline").
		SetLabelColor(tcell.ColorLightGray).
		SetBackgroundColor(tcell.ColorDarkSlateGrey).
		SetSelectedFunc(func() {
			p.State.Emit(events.EventQuit, nil)
		})

	formChoice := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 1, 1, false).
		AddItem(btnAccept, 12, 1, false).
		AddItem(nil, 3, 1, false).
		AddItem(btnDecline, 12, 1, false)

	p.BaseLayout().AddItem(viewTermsOfUse, 0, 3, false).
		AddItem(viewPrivacyPolicy, 0, 3, false).
		AddItem(formChoice, 1, 1, false)
}
