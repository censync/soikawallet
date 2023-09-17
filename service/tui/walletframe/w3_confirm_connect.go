package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
)

type pageW3ConfirmConnect struct {
	*BaseFrame
	*state.State
	// vars

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

	w3Chains := p.API().GetAllEvmChains(&dto.GetChainsDTO{
		OnlyW3: true,
	})
	labelSelectChains := tview.NewFlex().SetDirection(tview.FlexRow)

	for _, chain := range w3Chains {
		selectChainCheckbox := tview.NewCheckbox().
			SetLabel(chain.Name).
			SetChecked(true)
		labelSelectChains.AddItem(selectChainCheckbox, 1, 1, false)
	}

	layoutConnectionInfo := tview.NewFlex().SetDirection(tview.FlexRow)

	labelConnectionTitle := tview.NewTextView().SetLabel("Incoming new extension connection")

	labelConnectionInstanceId := tview.NewTextView().SetLabel(fmt.Sprintf("Instance ID: %s", connectionReq.InstanceId))

	labelConnectionOrigin := tview.NewTextView().SetLabel(fmt.Sprintf("Origin: %s", connectionReq.Origin))

	labelConnectionAddr := tview.NewTextView().SetLabel(fmt.Sprintf("Remote address: %s", connectionReq.RemoteAddr))

	layoutConnectionInfo.
		AddItem(labelConnectionTitle, 1, 1, false).
		AddItem(labelConnectionInstanceId, 1, 1, false).
		AddItem(labelConnectionOrigin, 1, 1, false).
		AddItem(labelConnectionAddr, 1, 1, false).
		AddItem(nil, 2, 1, false).
		AddItem(labelSelectChains, 0, 1, false)

	btnConnectionAccept := tview.NewButton(p.Tr().T("ui.button", "accept"))

	btnConnectionAccept.SetSelectedFunc(func() {
		p.EmitW3(event_bus.EventW3ConnAccepted, &dto.ResponseAcceptDTO{
			InstanceId: connectionReq.InstanceId,
			Chains:     w3Chains,
		})
		p.SwitchToPage(p.Pages().GetPrevious())
	})
	btnConnectionReject := tview.NewButton(p.Tr().T("ui.button", "reject"))

	btnConnectionReject.SetSelectedFunc(func() {
		p.EmitW3(event_bus.EventW3ConnRejected, &dto.ResponseRejectDTO{
			InstanceId: connectionReq.InstanceId,
			RemoteAddr: connectionReq.RemoteAddr,
		})
		p.SwitchToPage(p.Pages().GetPrevious())
	})

	layoutButtons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 3, 1, false).
		AddItem(btnConnectionAccept, 0, 1, false).
		AddItem(nil, 0, 1, false).
		AddItem(btnConnectionReject, 0, 1, false).
		AddItem(nil, 0, 2, false)

	layoutButtons.SetBorderPadding(0, 0, 10, 10)

	p.layout.
		AddItem(layoutConnectionInfo, 0, 1, false).
		AddItem(layoutButtons, 3, 1, false)
}

func (p *pageW3ConfirmConnect) FuncOnHide() {
	p.layout.Clear()
}
