package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/rivo/tview"
)

type pageW3RequestAccounts struct {
	*BaseFrame
	*state.State
}

func newPageW3RequestAccounts(state *state.State) *pageW3RequestAccounts {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3RequestAccounts{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageW3RequestAccounts) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Request address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	//connectionReq := p.Params()[0].(*dto.ConnectDTO)

}

func (p *pageW3RequestAccounts) FuncOnHide() {
	p.layout.Clear()
}
