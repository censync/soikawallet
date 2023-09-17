package walletframe

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/widgets/tabs"
	"github.com/censync/tview"
)

type pageSettings struct {
	*BaseFrame
	*state.State

	// vars
	layoutRPCList *tview.Flex
}

func newPageSettings(state *state.State) *pageSettings {
	layout := tview.NewFlex()

	return &pageSettings{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageSettings) FuncOnShow() {
	tabs := tabs.NewTabs().
		AddItem("Application", p.tabApp()).
		AddItem("Labels", p.tabLabels()).
		AddItem("RPC", p.tabNodes())
	p.layout.AddItem(tabs, 0, 1, false)

}

func (p *pageSettings) FuncOnHide() {
	p.layout.Clear()
}
