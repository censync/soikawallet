package walletframe

import (
	"github.com/censync/go-i18n"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/extpages"
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
}

func Init(events *handler.TBus, tr *i18n.Translator) *WalletFrame {

	frame := &WalletFrame{state: state.InitState(events, tr)}
	pages := frame.initPages()
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
	pages.SwitchToPage(pageNameSelectInitMode)
	return pages
}

func (f *WalletFrame) Layout() *tview.Flex {
	layoutMenu := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton("[yellow][F1] [white]QR Page     ", func() {
			f.state.SwitchToPage(pageNameQR)
			// t.events <- pageQR
		}).
		AddButton("[yellow][F5] [white]Addresses   ", func() {
			f.state.SwitchToPage(pageNameAddresses)
		}).
		AddButton("[yellow][F4] [white]Create      ", func() {
			if f.state.API() != nil {
				f.state.SwitchToPage(pageNameCreateWallets)
			}
		}).
		AddButton("[yellow][F6] [white]Transactions", func() {
			//if f.wallet != nil {
			f.state.SwitchToPage(pageNameTransaction)
			//}
		}).
		AddButton("[yellow][F8] [white]Node info   ", func() {
			f.state.SwitchToPage(pageNameNodeInfo)
		}).
		AddButton("[yellow][F3] [white]Settings    ", func() {
			f.state.SwitchToPage(pageNameSettings)
		}).
		AddButton("[yellow][F1] [white]About       ", func() {
			f.state.SwitchToPage(pageNameAbout)
		}).
		AddButton("[yellow][F10] [white]Quit       ", func() {
			f.state.Emit(handler.EventQuit, nil)

		})
	layoutMain := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(layoutMenu, 25, 1, false).
		AddItem(f.state.Pages(), 0, 1, true)

	layoutMenu.SetBorder(true).
		SetBorderColor(tcell.ColorDarkGrey)
	return layoutMain
}
