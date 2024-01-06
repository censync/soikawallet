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

package rpc

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageNodeInfo struct {
	*twoframes.BaseFrame
	*state.State

	// ui
	labelRPCInfo    *tview.TextView
	inputSelectNode *tview.DropDown

	// var
	selectedChain  *mhda.Chain
	availableNodes map[uint32]*resp.RPCInfo
	selectedNode   uint32
}

func NewPageNodeInfo(state *state.State) *pageNodeInfo {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(1, 0, 0, 0)

	return &pageNodeInfo{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageNodeInfo) FuncOnShow() {
	p.inputSelectNode = tview.NewDropDown().SetLabel("Select node")

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetFieldWidth(10).
		SetOptions(p.API().GetAllChainNames(), func(text string, index int) {
			p.selectedChain = p.API().GetChainByName(&dto.GetChainByNameDTO{
				ChainName: text,
			})
			p.availableNodes = p.API().AllRPC(&dto.GetRPCListByNetworkDTO{
				ChainKey: p.selectedChain.Key(),
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

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutResult, 0, 3, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageNodeInfo) actionUpdateNodesList() {
	nodesLabels := make([]string, 0)
	nodesIndex := map[int]uint32{}

	index := 0
	for nodeIndex, nodeInfo := range p.availableNodes {
		labelFormat := "#%d - %s"
		if nodeInfo.IsDefault {
			labelFormat = "[Default] " + labelFormat
		}
		nodesIndex[index] = nodeIndex
		nodesLabels = append(nodesLabels, fmt.Sprintf(labelFormat, nodeIndex, nodeInfo.Title))
		index++
	}

	p.inputSelectNode.SetOptions(nodesLabels, func(text string, index int) {
		p.selectedNode = nodesIndex[index]
	}).SetCurrentOption(0)
}

func (p *pageNodeInfo) actionUpdateInfo() {
	p.labelRPCInfo.Clear()
	receipt, err := p.API().GetRPCInfo(&dto.GetRPCInfoDTO{
		ChainKey:  p.selectedChain.Key(),
		NodeIndex: p.selectedNode,
	})
	if err == nil {
		str := ""
		for key, value := range receipt {
			str += fmt.Sprintf("[lightgray]%s: [darkcyan]%s\n", key, value)
		}
		p.labelRPCInfo.SetText(str)
	} else {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot get rpc info: %s", err))
	}

}
