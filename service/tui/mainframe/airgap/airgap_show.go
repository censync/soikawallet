// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package airgap

import (
	"fmt"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/qrview"
	"github.com/censync/twidget/twoframes"
	"github.com/gdamore/tcell/v2"
)

type pageAirGapShow struct {
	*twoframes.BaseFrame
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
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageAirGapShow) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(events.EventLogError, fmt.Sprintf("Params required"))
		p.SwitchToPage(p.Pages().GetPrevious())
		return
	}

	airGapData, ok := p.Params()[0].(*responses.AirGapMessage)

	if !ok {
		p.Emit(events.EventLogError, fmt.Sprintf("Incorrect params"))
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	p.Emit(events.EventLog, fmt.Sprintf("CHUNKS: %d", len(airGapData.Chunks)))

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(func() {

		})
	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 2, 1, false).
		AddItem(btnNext, 3, 1, false)

	redraw := func() {
		p.Emit(events.EventDrawForce, nil)
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
