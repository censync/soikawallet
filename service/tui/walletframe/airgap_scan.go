package walletframe

import (
	"fmt"
	"github.com/censync/go-zbar"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/rivo/tview"
)

type pageAirGapScan struct {
	*BaseFrame
	*state.State
}

func newPageAirGapScan(state *state.State) *pageAirGapScan {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	return &pageAirGapScan{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageAirGapScan) FuncOnShow() {
	scannerInstance := zbar.NewProcessor(1)

	scannerInstance.RequestSize(400, 400)
	scannerInstance.RequestInterface(2)
	defer scannerInstance.Destroy()

	scannerInstance.SetConfig(zbar.ZBAR_QRCODE, zbar.ZBAR_CFG_ENABLE, 1)

	if resultCode := scannerInstance.SetConfig(zbar.ZBAR_QRCODE, zbar.ZBAR_CFG_ENABLE, 1); resultCode != 0 {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot set config, code: %d", resultCode))
	}

	if resultCode := scannerInstance.Init("/dev/video0", 1); resultCode != 0 {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot qr scanner, code: %d", resultCode))
	}

	scannerInstance.SetDataHandler(p.handleScan)
	scannerInstance.SetActive(1)
	scannerInstance.SetVisible(1)

	scannerInstance.UserWait(-1)

}

func (p *pageAirGapScan) handleScan(img *zbar.Image) {
	s := img.FirstSymbol()
	if s != nil {
		p.Emit(event_bus.EventLogSuccess, fmt.Sprintf("QR: %s", s.GetData()))
	}
}

func (p *pageAirGapScan) FuncOnHide() {
}
