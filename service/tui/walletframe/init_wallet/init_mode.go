package init_wallet

import (
	"github.com/censync/soikawallet/service/tui/page"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageInitMode struct {
	*twidget.BaseFrame
	*state.State
}

func NewPageInitMode(state *state.State) *pageInitMode {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	layout.SetBorderPadding(5, 5, 5, 5)

	return &pageInitMode{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageInitMode) FuncOnShow() {
	btnModeAirGap := tview.NewButton("Use AirGap [[green]Recommended[white]]")
	btnModeAirGap.SetSelectedFunc(func() {
		p.SetWalletMode(state.ModeWithAirGap)
		p.SetStatus(state.StateInitAirGap)
		p.SwitchToPage(page.AirGapShow)
	})

	btnModeLocal := tview.NewButton("Do not use AirGap [[red]less secure[white]]")
	btnModeLocal.SetSelectedFunc(func() {
		p.SetWalletMode(state.ModeWithoutAirGap)
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(page.SelectInitWallet)
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
