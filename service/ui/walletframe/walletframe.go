package walletframe

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/extpages"
	"github.com/censync/soikawallet/service/ui/widgets/flexmenu"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageNameSelectInitMode   = `select_init_mode`
	pageNameSelectInitWallet = `select_init_wallet`
	pageNameInitMnemonic     = `init_mnemonic`
	pageNameRestoreMnemonic  = `restore_mnemonic`
	pageNameCreateWallets    = `wallets_create`
	pageNameAddresses        = `addresses`
	pageNameTransaction      = `transaction`
	pageNameOperationTx      = `operation_tx`
	pageNameTokenAdd         = `token_add`
	pageNameNodeInfo         = `node_info`
	pageNameSettings         = `settings`
	pageNameQR               = `qr`
	pageNameAgreement        = `agreement`
	pageNameAbout            = `about`
)

type IExtPage interface {
	Layout() *tview.Flex
	FuncOnShow()
	FuncOnHide()
	FuncOnDraw()
}

type BaseFrame struct {
	layout *tview.Flex
}

func (b *BaseFrame) Layout() *tview.Flex { return b.layout }

func (b *BaseFrame) FuncOnShow() {}

func (b *BaseFrame) FuncOnHide() {}

func (b *BaseFrame) FuncOnDraw() {}

type WalletFrame struct {
	state *state.State
	style *tview.Theme
}

func Init(events *handler.TBus, tr *i18n.Translator, style *tview.Theme) *WalletFrame {

	frame := &WalletFrame{state: state.InitState(events, tr), style: style}
	pages := frame.initPages()
	pages.SwitchToPage(pageNameAgreement)
	frame.state.SetPages(pages)
	return frame
}

func (f *WalletFrame) initPages() *extpages.ExtPages {
	prepared := map[string]IExtPage{
		pageNameSelectInitMode:   newPageInitMode(f.state),
		pageNameSelectInitWallet: newPageInitWallet(f.state),
		pageNameInitMnemonic:     newPageInitMnemonic(f.state),
		pageNameRestoreMnemonic:  newPageRestoreMnemonic(f.state),
		pageNameCreateWallets:    newPageCreateWallet(f.state),
		pageNameAddresses:        newPageAddresses(f.state),
		pageNameTransaction:      newPageTransactions(f.state),
		pageNameOperationTx:      newPageOperationTx(f.state),
		pageNameTokenAdd:         newPageTokenAdd(f.state),
		pageNameNodeInfo:         newPageNodeInfo(f.state),
		pageNameSettings:         newPageSettings(f.state),
		pageNameQR:               newPageQr(f.state),
		pageNameAgreement:        newPageAgreement(f.state),
		pageNameAbout:            newPageAbout(f.state),
	}
	pages := extpages.NewPages()

	for name, page := range prepared {
		pages.AddPage(extpages.NewPage(
			name,
			page.Layout(),
			true,
			false,
			page.FuncOnShow,
			page.FuncOnHide,
			page.FuncOnDraw,
		))
	}
	return pages
}

func (f *WalletFrame) Layout() *tview.Flex {
	layoutMenu := flexmenu.NewFlexMenu()

	layoutMenu.
		AddMenuItem("QR Dst", 0, func() { f.state.SwitchToPage(pageNameQR) }).
		AddMenuItem("Addresses", tcell.KeyF5, func() { f.state.SwitchToPage(pageNameAddresses) }).
		AddMenuItem("Create", tcell.KeyF2, func() { f.state.SwitchToPage(pageNameCreateWallets) }).
		AddMenuItem("Transactions", tcell.KeyF6, func() { f.state.SwitchToPage(pageNameTransaction) }).
		AddMenuItem("Node info", tcell.KeyF7, func() { f.state.SwitchToPage(pageNameNodeInfo) }).
		AddMenuItem("Settings", tcell.KeyF4, func() { f.state.SwitchToPage(pageNameSettings) }).
		AddMenuItem("About", tcell.KeyF1, func() { f.state.SwitchToPage(pageNameAbout) }).
		AddMenuItem("Quit", tcell.KeyF12, func() { f.state.Emit(handler.EventQuit, nil) })

	layoutMain := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(layoutMenu.Layout(), 25, 1, false).
		AddItem(f.state.Pages(), 0, 1, true)

	layoutMenu.SetBorder(true).
		SetBorderColor(tcell.ColorDarkGrey)
	return layoutMain
}
