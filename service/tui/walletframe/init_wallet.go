package walletframe

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/rivo/tview"
)

type pageInitWallet struct {
	*BaseFrame
	*state.State
}

func newPageInitWallet(state *state.State) *pageInitWallet {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageInitWallet{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageInitWallet) FuncOnShow() {
	btnWalletCreate := tview.NewButton(p.Tr().T("tui.button", "wallet_create"))

	btnWalletCreate.SetSelectedFunc(func() {
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(pageNameMnemonicInit)
	})
	btnWalletRestore := tview.NewButton(p.Tr().T("tui.button", "wallet_restore"))

	btnWalletRestore.SetSelectedFunc(func() {
		p.SetStatus(state.StateInitLocal)
		p.SwitchToPage(pageNameMnemonicRestore)
	})

	layoutButtons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(btnWalletCreate, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(btnWalletRestore, 0, 1, false).
		AddItem(nil, 0, 1, false)

	layoutButtons.SetBorderPadding(0, 0, 10, 10)

	p.layout.AddItem(layoutButtons, 3, 1, false)
}

func (p *pageInitWallet) FuncOnHide() {
	p.layout.Clear()
}
