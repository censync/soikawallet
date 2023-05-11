package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/qrview"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type pageQr struct {
	*BaseFrame
	*state.State

	// ui
	labelQR *qrview.QrView

	// vars
	paramAction uint8
}

func newPageQr(state *state.State) *pageQr {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	return &pageQr{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageQr) FuncOnShow() {
	/*if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			handler.EventLogError,
			fmt.Sprintf("Incorrect params"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}*/

	if p.paramAction == types.OpInitBootstrap {

	}
	chunks, err := p.API().ExportMeta()

	if err != nil {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get meta: %s", err))
	} else {
		p.Emit(event_bus.EventLog, fmt.Sprintf("CHUNKS: %d", len(chunks.Chunks)))
		redraw := func() {
			p.Emit(event_bus.EventDrawForce, nil)
		}
		p.labelQR = qrview.NewQrView(chunks.Chunks, 300, redraw)
		p.layout.AddItem(p.labelQR, 80, 1, false)
		p.labelQR.Start()
	}

}

func (p *pageQr) FuncOnHide() {
	p.labelQR.Stop()
	p.layout.Clear()
}
