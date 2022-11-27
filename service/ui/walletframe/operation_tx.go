package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
)

type pageOperationTx struct {
	*BaseFrame
	*state.State

	selectedAddr *responses.AddressResponse
}

func newPageOperationTx(state *state.State) *pageOperationTx {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageOperationTx{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageOperationTx) FuncOnShow() {
	for i := range p.Params() {
		p.Emit(handler.EventLogInfo, fmt.Sprintf("Params: %d = %s", i, p.Params()[i]))
	}
}
