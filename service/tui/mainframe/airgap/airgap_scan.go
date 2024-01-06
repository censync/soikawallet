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
	airgap "github.com/censync/go-airgap"
	"github.com/censync/go-zbar"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
	"github.com/gdamore/tcell/v2"
)

type pageAirGapScan struct {
	*twoframes.BaseFrame
	*state.State

	// vars
	scannerInstance *zbar.Processor
	chunks          *airgap.Chunks
	isScanStarted   bool
	chunksCounter   int
}

func NewPageAirGapScan(state *state.State) *pageAirGapScan {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	return &pageAirGapScan{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageAirGapScan) FuncOnShow() {
	btnNext := tview.NewButton(p.Tr().T("ui.button", "scan")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(func() {
			go p.actionScannerStart()
		})

	labelNotice := tview.NewTextView().
		SetText("Soika Vault")

	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 2, 1, false).
		AddItem(btnNext, 3, 1, false).
		AddItem(labelNotice, 3, 1, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutWizard, 30, 1, false).
		AddItem(nil, 0, 1, false)

	go p.actionScannerStart()
}

func (p *pageAirGapScan) actionScannerStart() {
	defer func() {
		if r := recover(); r != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Recovered in %s", r))
		}
	}()

	if p.isScanStarted {
		return
	}
	p.chunksCounter = 0
	p.isScanStarted = true

	p.chunks = airgap.NewChunks()
	p.scannerInstance = zbar.NewProcessor(1)

	// TODO: Add --nodbus flag

	p.scannerInstance.RequestSize(500, 500)
	//p.scannerInstance.RequestInterface(2)

	p.scannerInstance.SetConfig(zbar.ZBAR_QRCODE, zbar.ZBAR_CFG_ENABLE, 1)

	if resultCode := p.scannerInstance.SetConfig(zbar.ZBAR_QRCODE, zbar.ZBAR_CFG_ENABLE, 1); resultCode != 0 {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot set config, code: %d", resultCode))
	}

	if resultCode := p.scannerInstance.Init("/dev/video0", 1); resultCode != 0 {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot init_wallet camera, code: %d", resultCode))
		return
	}

	p.scannerInstance.SetDataHandler(p.handleScan)
	p.scannerInstance.SetActive(1)
	p.scannerInstance.SetVisible(1)

	p.scannerInstance.UserWait(-1)
}

func (p *pageAirGapScan) handleScan(img *zbar.Image) {
	s := img.FirstSymbol()
	if s != nil {
		wasAdded, err := p.chunks.ReadB64Chunk(s.GetData())
		if err != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Cannot read QR: %s", err))
		}

		if wasAdded {
			p.chunksCounter++
			p.Emit(events.EventLogInfo, fmt.Sprintf("Scanned [%d / %d]", p.chunksCounter, p.chunks.Count()))
		}

		if p.chunks.IsFilled() {
			p.Emit(events.EventLogSuccess, "Scan finished")
			p.actionScannerStop()

			p.actionProcessMessage()
		}
	}
}

func (p *pageAirGapScan) actionScannerStop() {
	p.isScanStarted = false
	if p.scannerInstance != nil {
		p.scannerInstance.SetActive(0)
		p.scannerInstance.SetVisible(0)
		p.scannerInstance.Destroy()
		p.scannerInstance = nil
	}
}

func (p *pageAirGapScan) actionProcessMessage() {
	result, err := p.API().ProcessAirGapMessage(&dto.AirGapMessageDTO{
		Data: p.chunks.Data(),
	})

	p.chunksCounter = 0
	p.chunks = nil

	if err == nil {
		p.Emit(events.EventLogSuccess, fmt.Sprintf("Operations scanned: %s", result))
	} else {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot process AirGap message: %s", err))
	}
}

func (p *pageAirGapScan) FuncOnHide() {
	p.actionScannerStop()
	p.BaseLayout().Clear()
}
