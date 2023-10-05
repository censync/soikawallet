package walletframe

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/page"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget/extpages"
	"github.com/censync/soikawallet/service/tui/twidget/flexmenu"
	"github.com/censync/soikawallet/service/tui/walletframe/about"
	"github.com/censync/soikawallet/service/tui/walletframe/addresses"
	"github.com/censync/soikawallet/service/tui/walletframe/airgap"
	"github.com/censync/soikawallet/service/tui/walletframe/create_addr"
	"github.com/censync/soikawallet/service/tui/walletframe/init_wallet"
	"github.com/censync/soikawallet/service/tui/walletframe/mnemonic"
	"github.com/censync/soikawallet/service/tui/walletframe/operation"
	"github.com/censync/soikawallet/service/tui/walletframe/rpc"
	"github.com/censync/soikawallet/service/tui/walletframe/settings"
	"github.com/censync/soikawallet/service/tui/walletframe/token"
	"github.com/censync/soikawallet/service/tui/walletframe/transaction"
	"github.com/censync/soikawallet/service/tui/walletframe/w3"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type IExtPage interface {
	Layout() *tview.Flex
	FuncOnShow()
	FuncOnHide()
	FuncOnDraw()
}

type WalletFrame struct {
	state *state.State
	style *tview.Theme
}

func Init(uiEvents, w3Events *event_bus.EventBus, tr *i18n.Translator, style *tview.Theme) *WalletFrame {

	frame := &WalletFrame{state: state.InitState(uiEvents, w3Events, tr), style: style}
	pages := frame.initPages()
	pages.SwitchToPage(page.Agreement)
	frame.state.SetPages(pages)
	return frame
}

func (f *WalletFrame) initPages() *extpages.ExtPages {
	prepared := map[string]IExtPage{
		page.SelectInitMode:   init_wallet.NewPageInitMode(f.state),
		page.SelectInitWallet: init_wallet.NewPageInitWallet(f.state),
		page.MnemonicInit:     mnemonic.NewPageInitMnemonic(f.state),
		page.MnemonicRestore:  mnemonic.NewPageRestoreMnemonic(f.state),
		page.CreateWallets:    create_addr.NewPageCreateWallet(f.state),
		page.Addresses:        addresses.NewPageAddresses(f.state),
		page.Transaction:      transaction.NewPageTransactions(f.state),
		page.OperationTx:      operation.NewPageOperationTx(f.state),
		page.TokenAdd:         token.NewPageTokenAdd(f.state),
		page.RPCInfo:          rpc.NewPageNodeInfo(f.state),
		page.Settings:         settings.NewPageSettings(f.state),
		page.AirGapShow:       airgap.NewPageAirGapShow(f.state),
		page.AirGapScan:       airgap.NewPageAirGapScan(f.state),
		page.Agreement:        about.NewPageAgreement(f.state),
		page.About:            about.NewPageAbout(f.state),
		// w3 connector
		page.W3ConfirmConnect:  w3.NewPageW3ConfirmConnect(f.state),
		page.W3RequestAccounts: w3.NewPageW3RequestAccounts(f.state),

		// internal
		page.W3Connections: w3.NewPageW3Connections(f.state),
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
	layoutMenu := flexmenu.NewFlexMenu(false)

	layoutMenu.
		// AddMenuItem("QR Show", 0, func() { f.state.SwitchToPage(AirGapShow) }).
		AddMenuItem("AirGap Scan", 0, func() { f.state.SwitchToPage(page.AirGapScan) }).
		AddMenuItem("Addresses", tcell.KeyF5, func() { f.state.SwitchToPage(page.Addresses) }).
		AddMenuItem("Create", tcell.KeyF2, func() { f.state.SwitchToPage(page.CreateWallets) }).
		AddMenuItem("Transactions", tcell.KeyF6, func() { f.state.SwitchToPage(page.Transaction) }).
		AddMenuItem("Node info", tcell.KeyF7, func() { f.state.SwitchToPage(page.RPCInfo) }).
		AddMenuItem("Settings", tcell.KeyF4, func() { f.state.SwitchToPage(page.Settings) }).
		AddMenuItem("W3 connections", tcell.KeyF3, func() { f.state.SwitchToPage(page.W3Connections) }).
		AddMenuItem("About", tcell.KeyF1, func() { f.state.SwitchToPage(page.About) }).
		AddMenuItem("Quit", tcell.KeyF12, func() { f.state.Emit(event_bus.EventQuit, nil) })

	layoutMain := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(layoutMenu, 25, 1, false).
		AddItem(f.state.Pages(), 0, 1, false)

	return layoutMain
}

func (f *WalletFrame) State() *state.State {
	return f.state
}
