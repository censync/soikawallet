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

package settings

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/config/chain"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/tview"
)

func (p *pageSettings) tabNodes() *tview.Flex {
	layoutShowRPCList := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	if layoutShowRPCList != nil {
		p.layoutRPCList = tview.NewFlex().SetDirection(tview.FlexRow)
	}

	p.actionUpdateRPCList()

	btnRefresh := tview.NewButton("refresh").
		SetSelectedFunc(p.actionUpdateRPCList)

	inputTitle := tview.NewInputField().SetLabel("Title")
	inputRPC := tview.NewInputField().SetLabel("RPC")

	btnAdd := tview.NewButton("Add").SetSelectedFunc(func() {
		err := p.API().AddRPC(&dto.AddRPCDTO{
			ChainKey: chain.EthereumChain.Key(), // TODO: Add selection chain key
			Title:    inputTitle.GetText(),
			Endpoint: inputRPC.GetText(),
		})
		if err != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Cannot add rpc \"%s\"", err))
		} else {
			p.Emit(events.EventLogInfo, fmt.Sprintf("Added rpc \"%s\"", inputTitle.GetText()))
			inputTitle.SetText(``)
			inputRPC.SetText(``)
			p.actionUpdateRPCList()
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
		AddItem(p.layoutRPCList, 0, 3, false).
		AddItem(layoutActions, 0, 1, false)

	return layoutShowRPCList
}

func (p *pageSettings) actionUpdateRPCList() {
	p.layoutRPCList.Clear()

	if p.API() != nil {
		rpcList := p.API().AllRPC(&dto.GetRPCListByNetworkDTO{
			ChainKey: chain.EthereumChain.Key(), // TODO: Add selection chain key
		})

		for index, rpc := range rpcList {
			linkedAccountCount := p.API().GetRPCLinkedAccountCount(&dto.GetRPCLinkedAccountCountDTO{
				ChainKey:  chain.EthereumChain.Key(), // TODO: Add selection chain key
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

			p.layoutRPCList.AddItem(rpcEntry, 5, 1, false)
		}
	}
}
