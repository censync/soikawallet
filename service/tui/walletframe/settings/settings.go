package settings

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/service/tui/twidget/tabs"
	"github.com/censync/tview"
)

type pageSettings struct {
	*twidget.BaseFrame
	*state.State

	// vars
	layoutRPCList *tview.Flex
}

func NewPageSettings(state *state.State) *pageSettings {
	layout := tview.NewFlex()

	return &pageSettings{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageSettings) FuncOnShow() {
	tabs := tabs.NewTabs().
		AddItem("Application", p.tabApp()).
		AddItem("Labels", p.tabLabels()).
		AddItem("RPC", p.tabNodes())
	p.BaseLayout().AddItem(tabs, 0, 1, false)

}
