package init_wallet

import (
	"github.com/censync/soikawallet/service/tui/page"
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
		p.SwitchToPage(page.MnemonicInit)
	})
	btnWalletRestore := tview.NewButton(p.Tr().T("ui.button", "wallet_restore"))

	btnWalletRestore.SetSelectedFunc(func() {
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(page.MnemonicRestore)
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
