package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
)

type pageW3ConfirmConnect struct {
	*BaseFrame
	*state.State
}

func newPageW3ConfirmConnect(state *state.State) *pageW3ConfirmConnect {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3ConfirmConnect{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageW3ConfirmConnect) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Request address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	connectionReq := p.Params()[0].(*dto.ConnectDTO)

	layoutConnectionInfo := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	labelConnectionInfo := tview.NewTextView().SetText(fmt.Sprintf(
		"Incomming new extension connection\n\n"+
			"Instance ID: %s\n"+
			"Origin: %s\n"+
			"Remote address: %s",
		connectionReq.InstanceId,
		connectionReq.Origin,
		connectionReq.RemoteAddr,
	))

	layoutConnectionInfo.
		AddItem(nil, 0, 1, false).
		AddItem(labelConnectionInfo, 0, 1, false).
		AddItem(nil, 0, 1, false)

	btnWalletCreate := tview.NewButton(p.Tr().T("ui.button", "accept"))

	btnWalletCreate.SetSelectedFunc(func() {
		p.EmitW3(event_bus.EventW3ConnAccepted, &dto.ResponseAcceptDTO{
			InstanceId: connectionReq.InstanceId,
		})
		p.SwitchToPage(p.Pages().GetPrevious())
	})
	btnWalletRestore := tview.NewButton(p.Tr().T("ui.button", "reject"))

	btnWalletRestore.SetSelectedFunc(func() {
		p.EmitW3(event_bus.EventW3ConnRejected, &dto.ResponseRejectDTO{
			InstanceId: connectionReq.InstanceId,
			RemoteAddr: connectionReq.RemoteAddr,
		})
		p.SwitchToPage(p.Pages().GetPrevious())
	})

	layoutButtons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 2, false).
		AddItem(btnWalletCreate, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(btnWalletRestore, 0, 1, false).
		AddItem(nil, 0, 2, false)

	layoutButtons.SetBorderPadding(0, 0, 10, 10)

	p.layout.
		AddItem(layoutConnectionInfo, 0, 1, false).
		AddItem(layoutButtons, 3, 1, false)
}

func (p *pageW3ConfirmConnect) FuncOnHide() {
	p.layout.Clear()
}
