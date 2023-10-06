package transaction

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/event_bus"
	"github.com/censync/tview"
)

type pageTransactions struct {
	*twidget.BaseFrame
	*state.State

	// ui
	labelTxReceipt  *tview.TextView
	inputSelectNode *tview.DropDown
	inputSelectedTx *tview.InputField

	selectedTx string

	// var
	selectedChain  *mhda.Chain
	availableNodes map[uint32]*types.RPC
	selectedNode   uint32
}

func NewPageTransactions(state *state.State) *pageTransactions {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(1, 0, 0, 0)

	return &pageTransactions{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageTransactions) FuncOnShow() {

	layoutReceipt := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layoutReceipt.SetBorder(true)

	searchForm := tview.NewForm().
		SetHorizontal(true)

	p.labelTxReceipt = tview.NewTextView()

	p.inputSelectNode = tview.NewDropDown().SetLabel("Select node")

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetFieldWidth(10).
		SetOptions(types.GetChainNames(), func(text string, index int) {
			p.selectedChain = types.GetChainByName(text)
			p.availableNodes = p.API().AllRPC(&dto.GetRPCListByNetworkDTO{
				ChainKey: p.selectedChain.Key(),
			})
			p.actionUpdateNodesList()
		}).
		SetCurrentOption(0)

	p.inputSelectedTx = tview.NewInputField().
		SetLabel("Hash").
		SetText(p.selectedTx).
		SetFieldWidth(70).
		SetChangedFunc(func(text string) {
			p.selectedTx = text
		})

	searchForm.
		AddFormItem(p.inputSelectedTx).
		AddButton("Clear", func() {
			p.inputSelectedTx.SetText("")
			p.selectedTx = ""
		}).
		AddButton("Get", func() {
			p.actionUpdateTxInfo()
		})
	optionsForm := tview.NewForm().
		SetHorizontal(true).
		AddFormItem(inputSelectNetwork).
		AddFormItem(p.inputSelectNode)

	layoutReceipt.AddItem(searchForm, 3, 1, false).
		AddItem(optionsForm, 3, 1, false).
		AddItem(p.labelTxReceipt, 0, 1, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutReceipt, 0, 4, false).
		AddItem(nil, 0, 1, false)

}

func (p *pageTransactions) actionUpdateNodesList() {
	nodesLabels := make([]string, 0)
	nodesIndex := map[int]uint32{}

	index := 0
	for nodeIndex, nodeInfo := range p.availableNodes {
		labelFormat := "#%d - %s"
		if nodeInfo.IsDefault() {
			labelFormat = "[Default] " + labelFormat
		}
		nodesIndex[index] = nodeIndex
		nodesLabels = append(nodesLabels, fmt.Sprintf(labelFormat, nodeIndex, nodeInfo.Title()))
		index++
	}

	p.inputSelectNode.SetOptions(nodesLabels, func(text string, index int) {
		p.selectedNode = nodesIndex[index]
	}).SetCurrentOption(0)
}

func (p *pageTransactions) actionUpdateTxInfo() {
	if p.API() != nil {
		p.labelTxReceipt.Clear()
		receipt, err := p.API().GetTxReceipt(&dto.GetTxReceiptDTO{
			ChainKey:  p.selectedChain.Key(),
			NodeIndex: p.selectedNode,
			Hash:      p.selectedTx,
		})
		if err == nil {
			str := ""
			for key, value := range receipt {
				str += fmt.Sprintf("%s: [darkcyan]%s\n", key, value)
			}
			p.labelTxReceipt.SetText(str)
		} else {
			p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get tx receipt: %s", err))
		}
	}
}

func (p *pageTransactions) FuncOnHide() {
	p.selectedTx = ``
	p.BaseLayout().Clear()
}
