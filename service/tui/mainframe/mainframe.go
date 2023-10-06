package mainframe

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/mainframe/about"
	"github.com/censync/soikawallet/service/tui/mainframe/addresses"
	"github.com/censync/soikawallet/service/tui/mainframe/airgap"
	"github.com/censync/soikawallet/service/tui/mainframe/create_addresses"
	"github.com/censync/soikawallet/service/tui/mainframe/init_wallet"
	"github.com/censync/soikawallet/service/tui/mainframe/mnemonic"
	"github.com/censync/soikawallet/service/tui/mainframe/operation"
	"github.com/censync/soikawallet/service/tui/mainframe/rpc"
	"github.com/censync/soikawallet/service/tui/mainframe/settings"
	"github.com/censync/soikawallet/service/tui/mainframe/token"
	"github.com/censync/soikawallet/service/tui/mainframe/transaction"
	"github.com/censync/soikawallet/service/tui/mainframe/w3"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget/extpages"
	"github.com/censync/soikawallet/service/tui/twidget/flexmenu"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type IExtPage interface {
	Layout() *tview.Flex
	FuncOnShow()
	FuncOnHide()
	FuncOnDraw()
}

type MainFrame struct {
	state *state.State
	style *tview.Theme
}

func Init(uiEvents, w3Events *events.EventBus, tr *i18n.Translator, style *tview.Theme) *MainFrame {

	frame := &MainFrame{state: state.InitState(uiEvents, w3Events, tr), style: style}
	appPages := frame.initPages()
	appPages.SwitchToPage(pages.Agreement)
	frame.state.SetPages(appPages)
	return frame
}

func (f *MainFrame) initPages() *extpages.ExtPages {
	prepared := map[string]IExtPage{
		pages.SelectInitMode:   init_wallet.NewPageInitMode(f.state),
		pages.SelectInitWallet: init_wallet.NewPageInitWallet(f.state),
		pages.MnemonicInit:     mnemonic.NewPageInitMnemonic(f.state),
		pages.MnemonicRestore:  mnemonic.NewPageRestoreMnemonic(f.state),
		pages.CreateWallets:    create_addresses.NewPageCreateWallet(f.state),
		pages.Addresses:        addresses.NewPageAddresses(f.state),
		pages.Transaction:      transaction.NewPageTransactions(f.state),
		pages.OperationTx:      operation.NewPageOperationTx(f.state),
		pages.TokenAdd:         token.NewPageTokenAdd(f.state),
		pages.RPCInfo:          rpc.NewPageNodeInfo(f.state),
		pages.Settings:         settings.NewPageSettings(f.state),
		pages.AirGapShow:       airgap.NewPageAirGapShow(f.state),
		pages.AirGapScan:       airgap.NewPageAirGapScan(f.state),
		pages.Agreement:        about.NewPageAgreement(f.state),
		pages.About:            about.NewPageAbout(f.state),
		// w3 connector
		pages.W3ConfirmConnect:  w3.NewPageW3ConfirmConnect(f.state),
		pages.W3RequestAccounts: w3.NewPageW3RequestAccounts(f.state),

		// internal
		pages.W3Connections: w3.NewPageW3Connections(f.state),
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

func (f *MainFrame) Layout() *tview.Flex {
	layoutMenu := flexmenu.NewFlexMenu(false)

	layoutMenu.
		// AddMenuItem("QR Show", 0, func() { f.state.SwitchToPage(AirGapShow) }).
		AddMenuItem("AirGap Scan", 0, func() { f.state.SwitchToPage(pages.AirGapScan) }).
		AddMenuItem("Addresses", tcell.KeyF5, func() { f.state.SwitchToPage(pages.Addresses) }).
		AddMenuItem("Create", tcell.KeyF2, func() { f.state.SwitchToPage(pages.CreateWallets) }).
		AddMenuItem("Transactions", tcell.KeyF6, func() { f.state.SwitchToPage(pages.Transaction) }).
		AddMenuItem("Node info", tcell.KeyF7, func() { f.state.SwitchToPage(pages.RPCInfo) }).
		AddMenuItem("Settings", tcell.KeyF4, func() { f.state.SwitchToPage(pages.Settings) }).
		AddMenuItem("W3 connections", tcell.KeyF3, func() { f.state.SwitchToPage(pages.W3Connections) }).
		AddMenuItem("About", tcell.KeyF1, func() { f.state.SwitchToPage(pages.About) }).
		AddMenuItem("Quit", tcell.KeyF12, func() { f.state.Emit(events.EventQuit, nil) })

	layoutMain := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(layoutMenu, 25, 1, false).
		AddItem(f.state.Pages(), 0, 1, false)

	return layoutMain
}

func (f *MainFrame) State() *state.State {
	return f.state
}
