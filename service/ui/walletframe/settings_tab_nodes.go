package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

func (p *pageSettings) tabNodes() *tview.Flex {
	layoutShowRPCList := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layoutRPCList := tview.NewFlex().SetDirection(tview.FlexRow)

	funcUpdate := func() {
		layoutRPCList.Clear()

		if p.API() != nil {
			rpcList := p.API().AllRPC(&dto.GetRPCListByCoinDTO{
				CoinType: uint32(types.Ethereum),
			})

			for index, rpc := range rpcList {
				linkedAccountCount := p.API().GetRPCLinkedAccountCount(&dto.GetRPCLinkedAccountCountDTO{
					CoinType:  uint32(types.Ethereum),
					NodeIndex: index,
				})
				btnEditEntry := tview.NewButton("Edit")
				btnRemoveEntry := tview.NewButton("Remove")

				nodeDescFormat := ""
				if rpc.IsDefault() {
					nodeDescFormat = fmt.Sprintf("#%d - %s\nRPC:%s\n[Default]", index, rpc.Title(), rpc.Endpoint())
				} else {
					nodeDescFormat = fmt.Sprintf("#%d - %s\nRPC:%s\nLinked accounts: %d", index, rpc.Title(), rpc.Endpoint(), linkedAccountCount)
				}

				rpcEntry := tview.NewFlex().
					SetDirection(tview.FlexColumn).
					AddItem(tview.NewTextView().SetText(nodeDescFormat), 0, 5, false).
					AddItem(btnEditEntry, 0, 1, false).
					AddItem(nil, 0, 1, false).
					AddItem(btnRemoveEntry, 0, 1, false)

				rpcEntry.SetBorder(true)

				layoutRPCList.AddItem(rpcEntry, 5, 1, false)
			}
		}
	}

	btnRefresh := tview.NewButton("refresh").SetSelectedFunc(funcUpdate)

	inputTitle := tview.NewInputField().SetLabel("Title")
	inputRPC := tview.NewInputField().SetLabel("RPC")

	btnAdd := tview.NewButton("Add").SetSelectedFunc(func() {
		err := p.API().AddRPC(&dto.AddRPCDTO{
			CoinType: 60,
			Title:    inputTitle.GetText(),
			Endpoint: inputRPC.GetText(),
		})
		if err != nil {
			p.Emit(handler.EventLogError, fmt.Sprintf("Cannot add rpc \"%s\"", err))
		} else {
			p.Emit(handler.EventLogInfo, fmt.Sprintf("Added rpc \"%s\"", inputTitle.GetText()))
			inputTitle.SetText(``)
			inputRPC.SetText(``)
			funcUpdate()
		}
	})

	layoutActions := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnRefresh, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(inputTitle, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(inputRPC, 1, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(btnAdd, 1, 1, false)

	layoutShowRPCList.
		//AddItem(layoutTableControls, 2, 1, false).
		AddItem(layoutRPCList, 0, 3, false).
		AddItem(layoutActions, 0, 1, false)

	return layoutShowRPCList
}
