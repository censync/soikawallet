//go:build !testnet

package walletframe

import (
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
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
		return
	}

	viewTermsOfUse := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(true).
		SetTextAlign(tview.AlignLeft).
		SetText(p.Tr().T("ui.label", "terms_of_use"))

	viewTermsOfUse.SetBorder(true).
		SetTitle(" Terms of use ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	viewPrivacyPolicy := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(false).
		SetTextAlign(tview.AlignLeft).
		SetText(p.Tr().T("ui.label", "privacy_policy"))

	viewPrivacyPolicy.SetBorder(true).
		SetTitle(" Privacy policy ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(0, 0, 1, 1)

	btnAccept := tview.NewButton("Accept").
		SetSelectedFunc(func() {
			p.SwitchToPage(pageNameSelectInitWallet)
		})

	btnDecline := tview.NewButton("Decline").
		SetLabelColor(tcell.ColorLightGray).
		SetBackgroundColor(tcell.ColorDarkSlateGrey).
		SetSelectedFunc(func() {
			p.State.Emit(event_bus.EventQuit, nil)
		})

	formChoice := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 1, 1, false).
		AddItem(btnAccept, 12, 1, false).
		AddItem(nil, 3, 1, false).
		AddItem(btnDecline, 12, 1, false)

	p.layout.AddItem(viewTermsOfUse, 0, 3, false).
		AddItem(viewPrivacyPolicy, 0, 3, false).
		AddItem(formChoice, 1, 1, false)
}

func (p *pageAgreement) FuncOnHide() {
	p.layout.Clear()
}
