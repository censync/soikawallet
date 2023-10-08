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

package addresses

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/service/tui/twidget/spinner"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

const (
	addrNodeLevelChain   = 1
	addrNodeLevelAccount = 2
	addrNodeLevelAddr    = 3
)

type pageAddresses struct {
	*twidget.BaseFrame
	*state.State

	// ui

	layoutAddressesTree *tview.TreeView
	layoutDetails       *tview.Flex
	addrTree            *tview.TreeNode

	// var
	selectedAddress *resp.AddressResponse
	selectedAccount *resp.AccountResponse

	isUpdating     bool
	balanceSpinner *spinner.Spinner
}

func NewPageAddresses(state *state.State) *pageAddresses {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageAddresses{
		State:               state,
		BaseFrame:           twidget.NewBaseFrame(layout),
		layoutAddressesTree: tview.NewTreeView(),
		balanceSpinner:      spinner.NewSpinner(spinner.SpinThree, 180),
	}
}

func (p *pageAddresses) Layout() *tview.Flex {

	p.addrTree = tview.NewTreeNode("wallets")
	p.layoutAddressesTree.SetRoot(p.addrTree).SetTopLevel(1)

	// double click for address operations
	p.layoutAddressesTree.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if p.layoutAddressesTree.InRect(event.Position()) {
			if action == tview.MouseLeftDoubleClick && p.selectedAddress != nil {
				if p.layoutAddressesTree.GetCurrentNode().GetLevel() == addrNodeLevelAddr {
					p.SwitchToPage(pages.OperationTx, p.selectedAddress)
				}
				return action, nil
			}
		}
		return action, event
	})

	p.layoutAddressesTree.SetBorder(true)

	p.layoutAddressesTree.SetSelectedFunc(func(node *tview.TreeNode) {
		// p.clearLayoutSelected()
		p.actionUpdateFrameDetails()

		reference := node.GetReference()
		p.selectedAddress = nil
		p.selectedAccount = nil

		switch reference.(type) {
		case *addrNodeViewEntry:
			p.selectedAddress = reference.(*addrNodeViewEntry).addr
		case *accountNodeViewEntry:
			p.selectedAccount = reference.(*accountNodeViewEntry).account
		}

		p.actionUpdateFrameDetails()
	})

	p.layoutDetails = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.BaseLayout().
		AddItem(p.layoutAddressesTree, 0, 1, false).
		AddItem(p.layoutDetails, 50, 1, false)

	return p.BaseFrame.Layout()
}

func (p *pageAddresses) actionUpdateFrameDetails() {
	if p.BaseLayout().GetItemCount() != 2 {
		return
	}
	detailsFrame := tview.NewFlex()
	if p.selectedAddress != nil {
		frame := newFrameAddressesDetailsAddr(p.State, p.selectedAddress)
		detailsFrame = frame.Layout()
	} else if p.selectedAccount != nil {
		frame := newAddressesFrameDetailsAccount(p.State, p.selectedAccount)
		detailsFrame = frame.Layout()
	} else {
		frame := newFrameAddressesDetailsEmpty(p.State)
		detailsFrame = frame.Layout()
	}
	item := p.BaseLayout().GetItem(1)
	p.BaseLayout().RemoveItem(item)
	p.BaseLayout().AddItem(detailsFrame, 50, 1, false)
}

func (p *pageAddresses) FuncOnShow() {
	go p.actionUpdateAddresses()
}

func (p *pageAddresses) FuncOnHide() {
	// DO NOT REMOVE
}
