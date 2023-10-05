package airgap

import (
	"fmt"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/service/tui/twidget/qrview"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type pageAirGapShow struct {
	*twidget.BaseFrame
	*state.State

	// ui
	labelQR *qrview.QrView

	// vars
	paramAction uint8
}

func NewPageAirGapShow(state *state.State) *pageAirGapShow {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	return &pageAirGapShow{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageAirGapShow) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Params required"))
		p.SwitchToPage(p.Pages().GetPrevious())
		return
	}

	airGapData, ok := p.Params()[0].(*responses.AirGapMessage)

	if !ok {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Incorrect params"))
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	p.Emit(event_bus.EventLog, fmt.Sprintf("CHUNKS: %d", len(airGapData.Chunks)))

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(func() {

		})
	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 2, 1, false).
		AddItem(btnNext, 3, 1, false)

	redraw := func() {
		p.Emit(event_bus.EventDrawForce, nil)
	}

	p.labelQR = qrview.NewQrView(airGapData.Chunks, 300, redraw)

	p.BaseLayout().AddItem(p.labelQR, 80, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(layoutWizard, 30, 1, false).
		AddItem(nil, 0, 1, false)

	p.labelQR.Start()
}

func (p *pageAirGapShow) FuncOnHide() {
	if p.labelQR != nil {
		p.labelQR.Stop()
	}
	p.BaseLayout().Clear()
}
