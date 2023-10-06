package create_addresses

import (
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/service/tui/twidget/tabs"
	"github.com/censync/tview"
)

type pageCreateAddr struct {
	*twidget.BaseFrame
	*state.State

	// ui
	// wizard
	layoutAddrEntriesForm     *tview.Flex
	inputSelectDerivationType *tview.DropDown
	inputSelectDerivationPath *tview.TextView
	inputUseHardenedAddresses *tview.Checkbox
	inputAccountIndex         *tview.InputField
	inputAddrIndex            *tview.InputField

	// bulk
	inputDerivationPaths *tview.TextArea

	// var
	// wizard
	selectedChain          *mhda.Chain
	selectedDerivationType mhda.DerivationType
	selectedDerivationPath mhda.DerivationPath
	selectedCharge         uint8
	selectedUseHardened    bool
	addrPoolGap            int
	accountStartIndex      int
	addrStartIndex         int
	// bulk
	//rxAddressPath *regexp.Regexp
}

func NewPageCreateWallet(state *state.State) *pageCreateAddr {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageCreateAddr{
		State:       state,
		BaseFrame:   twidget.NewBaseFrame(layout),
		addrPoolGap: defaultAddrPoolGap,
		//rxAddressPath: regexp.MustCompile(`(m/44[Hh']/[0-9]+[Hh']/[0-9]+[Hh']/[0|1]/[0-9]+[Hh']*)`),
	}
}

func (p *pageCreateAddr) FuncOnShow() {
	tabs := tabs.NewTabs().
		AddItem(p.Tr().T("ui.tab", "wizard"), p.tabWizard()).
		AddItem(p.Tr().T("ui.tab", "bulk"), p.tabBulk())
	p.BaseLayout().AddItem(tabs, 0, 1, false)
}
