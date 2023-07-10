package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type pageNodeInfo struct {
	*BaseFrame
	*state.State

	// ui
	labelRPCInfo    *tview.TextView
	inputSelectNode *tview.DropDown

	// var
	selectedChain  types.CoinType
	availableNodes map[uint32]*types.RPC
	selectedNode   uint32
}

func newPageNodeInfo(state *state.State) *pageNodeInfo {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(1, 0, 0, 0)

	return &pageNodeInfo{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageNodeInfo) FuncOnShow() {
	p.inputSelectNode = tview.NewDropDown().SetLabel("Select node")

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetFieldWidth(10).
		SetOptions(types.GetCoinNames(), func(text string, index int) {
			p.selectedChain = types.GetCoinByName(text)
			p.availableNodes = p.API().AllRPC(&dto.GetRPCListByCoinDTO{
				CoinType: uint32(p.selectedChain),
			})
			p.actionUpdateNodesList()
		}).
		SetCurrentOption(0)

	layoutResult := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layoutResult.SetBorder(true)

	searchForm := tview.NewForm().
		SetHorizontal(true).
		AddFormItem(inputSelectNetwork).
		AddFormItem(p.inputSelectNode).
		AddButton("Get", func() {
			p.actionUpdateInfo()
		})

	p.labelRPCInfo = tview.NewTextView().
		SetDynamicColors(true)

	layoutResult.AddItem(searchForm, 3, 1, false).
		AddItem(p.labelRPCInfo, 0, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutResult, 0, 3, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageNodeInfo) actionUpdateNodesList() {
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

func (p *pageNodeInfo) actionUpdateInfo() {
	p.labelRPCInfo.Clear()
	receipt, err := p.API().GetRPCInfo(&dto.GetRPCInfoDTO{
		CoinType:  uint32(p.selectedChain),
		NodeIndex: p.selectedNode,
	})
	if err == nil {
		str := ""
		for key, value := range receipt {
			str += fmt.Sprintf("[lightgray]%s: [darkcyan]%s\n", key, value)
		}
		p.labelRPCInfo.SetText(str)
	} else {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get rpc info: %s", err))
	}

}

func (p *pageNodeInfo) FuncOnHide() {
	p.layout.Clear()
}
