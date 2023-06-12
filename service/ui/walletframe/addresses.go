package walletframe

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/spinner"
	"github.com/censync/soikawallet/types"
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
	selectedCoin    types.CoinType
	selectedAccount types.AccountIndex
	isUpdating      bool
	balanceSpinner  *spinner.Spinner
}

type accountNodeViewEntry struct {
	coinType     types.CoinType
	accountIndex types.AccountIndex
}

type addrNodeViewEntry struct {
	addr     *resp.AddressResponse
	balances *int // *resp.AddressTokensBalanceListResponse
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
		p.selectedCoin = 0
		p.selectedAccount = 0

		if reference != nil {
			if addressEntry, ok := reference.(*addrNodeViewEntry); ok {
				p.Emit(event_bus.EventLog, "Addr selected")
				p.selectedAddress = addressEntry.addr
			} else if accountEntry, ok := reference.(*accountNodeViewEntry); ok {
				p.Emit(event_bus.EventLog, "Account selected")
				p.selectedCoin = accountEntry.coinType
				p.selectedAccount = accountEntry.accountIndex
			}
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
	} else if p.selectedCoin != 0 {
		p.Emit(event_bus.EventLog, "Account frame")
		frame := newFrameDetailsAccount(p.State, p.selectedCoin, p.selectedAccount)
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
