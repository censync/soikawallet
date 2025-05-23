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

package op_tx

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageOperationTx struct {
	*twoframes.BaseFrame
	*state.State

	paramSelectedAddr *responses.AddressResponse

	// ui
	layoutTokensTreeView *tview.TreeView
	layoutFrameOperation *tview.Flex

	// vars
}

func NewPageOperationTx(state *state.State) *pageOperationTx {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageOperationTx{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageOperationTx) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Sender address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.paramSelectedAddr = p.Params()[0].(*responses.AddressResponse)

	layoutOperation := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	layoutOperation.SetBorder(true)

	p.layoutTokensTreeView = tview.NewTreeView()

	frame := newFrameOperationWizard(p.State, p.paramSelectedAddr)
	p.layoutFrameOperation = frame.Layout()

	layoutOperation.AddItem(p.layoutTokensTreeView, 0, 1, false)
	layoutOperation.AddItem(p.layoutFrameOperation, 0, 2, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutOperation, 0, 4, false).
		AddItem(nil, 0, 1, false)

	go p.actionUpdateTokens()
}

func (p *pageOperationTx) actionUpdateTokens() {
	nodeTokens := tview.NewTreeNode("tokens")
	p.layoutTokensTreeView.SetRoot(nodeTokens).
		SetTopLevel(1)
	p.layoutTokensTreeView.SetBorder(true)

	balances, err := p.API().GetTokensBalancesByAddress(&dto.GetAddressTokensByPathDTO{
		MhdaPath: p.paramSelectedAddr.Path,
	})

	if err != nil {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Cannot get data for %s: %s", p.paramSelectedAddr.Path, err),
		)
	}

	for key, value := range balances {
		tokenNode := tview.NewTreeNode(fmt.Sprintf("$%s - %f", key, value))
		nodeTokens.AddChild(tokenNode)
	}

	p.Emit(events.EventDrawForce, nil)
}

func (p *pageOperationTx) FuncOnHide() {
	p.paramSelectedAddr = nil
	p.BaseLayout().Clear()
}
