package walletframe

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/spinner"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	addrNodeLevelChain   = 1
	addrNodeLevelAccount = 2
	addrNodeLevelAddr    = 3
)

type pageAddresses struct {
	*BaseFrame
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

func newPageAddresses(state *state.State) *pageAddresses {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageAddresses{
		State:               state,
		BaseFrame:           &BaseFrame{layout: layout},
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
					p.SwitchToPage(pageNameOperationTx, p.selectedAddress)
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

	p.layout.
		AddItem(p.layoutAddressesTree, 0, 1, false).
		AddItem(p.layoutDetails, 50, 1, false)

	return p.layout
}

func (p *pageAddresses) actionUpdateFrameDetails() {
	if p.layout.GetItemCount() != 2 {
		return
	}
	detailsFrame := tview.NewFlex()
	if p.selectedAddress != nil {
		frame := newFrameDetailsAddr(p.State, p.selectedAddress)
		detailsFrame = frame.Layout()
	} else if p.selectedAccount != nil {
		frame := newFrameDetailsAccount(p.State, p.selectedAccount)
		detailsFrame = frame.Layout()
	} else {
		frame := newFrameDetailsEmpty(p.State)
		detailsFrame = frame.Layout()
	}
	item := p.layout.GetItem(1)
	p.layout.RemoveItem(item)
	p.layout.AddItem(detailsFrame, 50, 1, false)
}

func (p *pageAddresses) FuncOnShow() {
	go p.actionUpdateAddresses()
}
