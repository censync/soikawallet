package walletframe

import (
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/qrview"
	"github.com/rivo/tview"
)

type pageQr struct {
	*BaseFrame
	*state.State

	// ui
	labelQR *tview.TextView
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
	addr := "0x5b4c4A2aE5721c9A2b994f5b4150C17EBB2f89E4"

	redraw := func() {
		p.Emit(handler.EventDrawForce, nil)
	}
	p.labelQR = qrview.ShowAnimation(&addr, 100, redraw)
	p.layout.AddItem(p.labelQR, 80, 1, false)
}

func (p *pageQr) FuncOnHide() {
	p.layout.Clear()
}
