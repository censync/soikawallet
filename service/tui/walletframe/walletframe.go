package walletframe

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/widgets/extpages"
	"github.com/censync/soikawallet/service/tui/widgets/flexmenu"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pageNameSelectInitMode   = `select_init_mode`
	pageNameSelectInitWallet = `select_init_wallet`
	pageNameMnemonicInit     = `mnemonic_init`
	pageNameMnemonicRestore  = `mnemonic_restore`
	pageNameCreateWallets    = `wallets_create`
	pageNameAddresses        = `addresses`
	pageNameTransaction      = `transaction`
	pageNameOperationTx      = `operation_tx`
	pageNameTokenAdd         = `token_add`
	pageNameNodeInfo         = `node_info`
	pageNameSettings         = `settings`
	pageNameAirGapShow       = `airgap_show`
	pageNameAirGapScan       = `airgap_scan`
	pageNameAgreement        = `agreement`
	pageNameAbout            = `about`

	// connector
	pageNameW3ConfirmConnect  = "w3_confirm_connect"
	pageNameW3RequestAccounts = "w3_request_accounts"
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

func Init(uiEvents, w3Events *event_bus.EventBus, tr *i18n.Translator, style *tview.Theme) *WalletFrame {

	frame := &WalletFrame{state: state.InitState(uiEvents, w3Events, tr), style: style}
	pages := frame.initPages()
	pages.SwitchToPage(pageNameAgreement)
	frame.state.SetPages(pages)
	return frame
}

func (f *WalletFrame) initPages() *extpages.ExtPages {
	prepared := map[string]IExtPage{
		pageNameSelectInitMode:   newPageInitMode(f.state),
		pageNameSelectInitWallet: newPageInitWallet(f.state),
		pageNameMnemonicInit:     newPageInitMnemonic(f.state),
		pageNameMnemonicRestore:  newPageRestoreMnemonic(f.state),
		pageNameCreateWallets:    newPageCreateWallet(f.state),
		pageNameAddresses:        newPageAddresses(f.state),
		pageNameTransaction:      newPageTransactions(f.state),
		pageNameOperationTx:      newPageOperationTx(f.state),
		pageNameTokenAdd:         newPageTokenAdd(f.state),
		pageNameNodeInfo:         newPageNodeInfo(f.state),
		pageNameSettings:         newPageSettings(f.state),
		pageNameAirGapShow:       newPageAirGapShow(f.state),
		pageNameAirGapScan:       newPageAirGapScan(f.state),
		pageNameAgreement:        newPageAgreement(f.state),
		pageNameAbout:            newPageAbout(f.state),
		// w3 connector
		pageNameW3ConfirmConnect:  newPageW3ConfirmConnect(f.state),
		pageNameW3RequestAccounts: newPageW3RequestAccounts(f.state),
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
		// AddMenuItem("QR Show", 0, func() { f.state.SwitchToPage(pageNameAirGapShow) }).
		AddMenuItem("AirGap Scan", 0, func() { f.state.SwitchToPage(pageNameAirGapScan) }).
		AddMenuItem("Addresses", tcell.KeyF5, func() { f.state.SwitchToPage(pageNameAddresses) }).
		AddMenuItem("Create", tcell.KeyF2, func() { f.state.SwitchToPage(pageNameCreateWallets) }).
		AddMenuItem("Transactions", tcell.KeyF6, func() { f.state.SwitchToPage(pageNameTransaction) }).
		AddMenuItem("Node info", tcell.KeyF7, func() { f.state.SwitchToPage(pageNameNodeInfo) }).
		AddMenuItem("Settings", tcell.KeyF4, func() { f.state.SwitchToPage(pageNameSettings) }).
		AddMenuItem("About", tcell.KeyF1, func() { f.state.SwitchToPage(pageNameAbout) }).
		AddMenuItem("Quit", tcell.KeyF12, func() { f.state.Emit(event_bus.EventQuit, nil) })

	layoutMain := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(layoutMenu, 25, 1, false).
		AddItem(f.state.Pages(), 0, 1, false)

	/*layoutMenu.SetBorder(true).
	SetBorderColor(tcell.ColorDarkGrey) */
	return layoutMain
}

func (f *WalletFrame) State() *state.State {
	return f.state
}
