package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type pageNodeInfo struct {
	*BaseFrame
	*state.State

	// ui
	labelRPCInfo *tview.TextView

	// var
	selectedChain types.CoinType
	selectedNode  uint32
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
	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetFieldWidth(10).
		SetOptions(types.GetCoinNames(), func(text string, index int) {
			p.selectedChain = types.GetCoinByName(text)
		}).
		SetCurrentOption(0)

	layoutReceipt := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layoutReceipt.SetBorder(true)

	searchForm := tview.NewForm().
		SetHorizontal(true).
		AddFormItem(inputSelectNetwork).
		AddButton("Get", func() {
			p.updateInfo()
		})

	p.labelRPCInfo = tview.NewTextView().
		SetDynamicColors(true)

	layoutReceipt.AddItem(searchForm, 3, 1, false).
		AddItem(p.labelRPCInfo, 0, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutReceipt, 0, 3, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageNodeInfo) updateInfo() {
	p.labelRPCInfo.Clear()
	receipt, err := p.API().GetRPCInfo(&dto.GetRPCInfoDTO{
		CoinType:  uint32(p.selectedChain),
		NodeIndex: 0,
	})
	if err == nil {
		str := ""
		for key, value := range receipt {
			str += fmt.Sprintf("[lightgray]%s: [darkcyan]%s\n", key, value)
		}
		p.labelRPCInfo.SetText(str)
	} else {
		p.Emit(handler.EventLogError, fmt.Sprintf("Cannot get rpc info: %s", err))
	}

}

func (p *pageNodeInfo) FuncOnHide() {
	p.layout.Clear()
}
