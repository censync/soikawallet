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
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package w3

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageW3ConfirmConnect struct {
	*twoframes.BaseFrame
	*state.State
}

func NewPageW3ConfirmConnect(state *state.State) *pageW3ConfirmConnect {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3ConfirmConnect{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageW3ConfirmConnect) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Request address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	connectionReq := p.Params()[0].(*dto.ConnectDTO)

	w3Chains := p.API().GetAllEvmW3Chains()
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
		p.EmitW3(events.EventW3ConnAccepted, &dto.ResponseAcceptDTO{
			InstanceId: connectionReq.InstanceId,
			Chains:     w3Chains,
		})
		p.SwitchToPage(p.Pages().GetPrevious())
	})
	btnConnectionReject := tview.NewButton(p.Tr().T("ui.button", "reject"))

	btnConnectionReject.SetSelectedFunc(func() {
		p.EmitW3(events.EventW3ConnRejected, &dto.ResponseRejectDTO{
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

	p.BaseLayout().
		AddItem(layoutConnectionInfo, 0, 1, false).
		AddItem(layoutButtons, 3, 1, false)
}
