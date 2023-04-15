package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/qrview"
	"github.com/rivo/tview"
)

type pageQr struct {
	*BaseFrame
	*state.State

	// ui
	labelQR *qrview.QrView
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
	chunks, err := p.API().ExportMeta()

	if err != nil {
		p.Emit(handler.EventLogError, fmt.Sprintf("Cannot get meta: %s", err))
	} else {
		p.Emit(handler.EventLog, fmt.Sprintf("CHUNKS: %d", len(chunks.Chunks)))
		redraw := func() {
			p.Emit(handler.EventDrawForce, nil)
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
