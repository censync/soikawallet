//go:build testnet

package walletframe

import (
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
)

type pageAgreement struct {
	*BaseFrame
	*state.State

	isAgreementAccepted bool
}

func newPageAgreement(state *state.State) *pageAgreement {
	var isAgreementAccepted bool

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(2, 2, 3, 3)

	envAgreementAccepted, ok := os.LookupEnv("SOIKAWALLET_AGREEMENT_ACCEPTED")

	if ok && strings.ToLower(envAgreementAccepted) == "true" {
		isAgreementAccepted = true
	}

	return &pageAgreement{
		BaseFrame:           &BaseFrame{layout: layout},
		State:               state,
		isAgreementAccepted: isAgreementAccepted,
	}
}

func (p *pageAgreement) FuncOnShow() {
	if p.isAgreementAccepted {
		p.SwitchToPage(pageNameSelectInitWallet)
	}

	viewTermsOfUse := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(false).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorRed).
		SetText(p.Tr().T("ui.label", "terms_of_use_testnet"))

	viewTermsOfUse.SetBorder(true)

	viewPrivacyPolicy := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(false).
		SetTextAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorRed).
		SetText(p.Tr().T("ui.label", "privacy_policy_testnet"))

	viewPrivacyPolicy.SetBorder(true)

	formChoice := tview.NewForm().
		SetHorizontal(true).
		SetButtonsAlign(tview.AlignCenter).
		AddButton("Accept", func() {
			p.SwitchToPage(pageNameSelectInitWallet)
		}).
		AddButton("Decline", func() {
			p.State.Emit(handler.EventQuit, nil)
		})
	p.layout.AddItem(viewTermsOfUse, 0, 1, false).
		AddItem(viewPrivacyPolicy, 0, 1, false).
		AddItem(formChoice, 0, 1, false)
}

func (p *pageAgreement) FuncOnHide() {
	p.layout.Clear()
}
