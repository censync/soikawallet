package walletframe

import (
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/tabs"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
	"regexp"
)

const addressPoolGap = 5

type pageCreateWallet struct {
	*BaseFrame
	*state.State

	// ui
	// wizard
	layoutCreateWalletsForm *tview.Flex
	// bulk
	inputDerivationPaths *tview.TextArea

	// var
	// wizard
	selectedChain       types.CoinType
	selectedCharge      uint8
	selectedUseHardened bool

	// bulk
	rxAddressPath *regexp.Regexp
}

func newPageCreateWallet(state *state.State) *pageCreateWallet {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageCreateWallet{
		State:         state,
		BaseFrame:     &BaseFrame{layout: layout},
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
