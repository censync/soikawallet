package operation

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageOperationTx struct {
	*twidget.BaseFrame
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
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageOperationTx) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			event_bus.EventLogError,
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
			event_bus.EventLogError,
			fmt.Sprintf("Cannot get data for %s: %s", p.paramSelectedAddr.Path, err),
		)
	}

	for key, value := range balances {
		tokenNode := tview.NewTreeNode(fmt.Sprintf("$%s - %f", key, value))
		nodeTokens.AddChild(tokenNode)
	}

	p.Emit(event_bus.EventDrawForce, nil)
}

func (p *pageOperationTx) FuncOnHide() {
	p.paramSelectedAddr = nil
	p.BaseLayout().Clear()
}
