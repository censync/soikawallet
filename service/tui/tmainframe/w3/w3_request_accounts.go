package w3

import (
	"fmt"
	"github.com/censync/soikawallet/service/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageW3RequestAccounts struct {
	*twidget.BaseFrame
	*state.State
}

func NewPageW3RequestAccounts(state *state.State) *pageW3RequestAccounts {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3RequestAccounts{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
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
