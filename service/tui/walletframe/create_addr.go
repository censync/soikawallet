package walletframe

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/widgets/tabs"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
	"regexp"
)

type pageCreateWallet struct {
	*BaseFrame
	*state.State

	// tui
	// wizard
	layoutAddrEntriesForm *tview.Flex
	// bulk
	inputDerivationPaths *tview.TextArea

	// var
	// wizard
	selectedChain       types.CoinType
	selectedCharge      uint8
	selectedUseHardened bool
	addrPoolGap         int
	accountStartIndex   int
	addrStartIndex      int
	// bulk
	rxAddressPath *regexp.Regexp
}

func newPageCreateWallet(state *state.State) *pageCreateWallet {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageCreateWallet{
		State:         state,
		BaseFrame:     &BaseFrame{layout: layout},
		addrPoolGap:   defaultAddrPoolGap,
		rxAddressPath: regexp.MustCompile(`(m/44[Hh']/[0-9]+[Hh']/[0-9]+[Hh']/[0|1]/[0-9]+[Hh']*)`),
	}
}

func (p *pageCreateWallet) FuncOnShow() {
	tabs := tabs.NewTabs().
		AddItem("Wizard", p.tabWizard()).
		AddItem("Bulk", p.tabBulk())
	p.layout.AddItem(tabs, 0, 1, false)
}

func (p *pageCreateWallet) FuncOnHide() {
	p.layout.Clear()
}
