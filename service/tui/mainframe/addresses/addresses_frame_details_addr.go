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

package addresses

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget/qrview"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type frameAddressesDetailsAddr struct {
	layout *tview.Flex
	*state.State

	// ui
	layoutAddrSelected *tview.Flex
	labelQR            *tview.TextView

	// vars
	selectedAddress *resp.AddressResponse
	isQrIsShown     bool
}

func newFrameAddressesDetailsAddr(state *state.State, selectedAddress *resp.AddressResponse) *frameAddressesDetailsAddr {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	return &frameAddressesDetailsAddr{
		State:           state,
		layout:          layout,
		selectedAddress: selectedAddress,
	}
}

func (f *frameAddressesDetailsAddr) Layout() *tview.Flex {

	f.layoutAddrSelected = tview.NewFlex().
		SetDirection(tview.FlexRow)

	//SetBorderPadding(0, 0, 0, 0)
	f.layoutAddrSelected.SetTitle(" Send from ").
		SetBorder(true).
		SetBorderColor(tcell.ColorDimGrey).
		SetBorderPadding(0, 0, 1, 0)

	f.labelQR = tview.NewTextView().
		SetWordWrap(false).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	// Selected address label
	// clipboard.CopyToClipboard(inputMnemonic.GetText())

	// TODO: Optimize to selected addr
	pathTitle := ""
	addr, err := mhda.ParseNSS(f.selectedAddress.Path)
	if err == nil {
		pathTitle = types.GetNetworkNameByKey(addr.Chain().Key()) + " " + addr.DerivationPath().String()
	}

	viewSelectedPath := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorLightBlue).
		SetText(pathTitle)

	viewSelectedAddr := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorBlue).
		SetText(f.selectedAddress.Address)

	btnSelectedAddrCopy := tview.NewButton("copy").SetSelectedFunc(func() {
		if f.selectedAddress != nil {
			err := clipboard.CopyToClipboard(f.selectedAddress.Address)
			if err != nil {
				f.Emit(events.EventLogError, fmt.Sprintf("Cannot copy to clipboard: %s", err))
			} else {
				f.Emit(events.EventLogSuccess, "Address copied to clipboard")
			}
		}
	})

	layoutSelectedAddrOptions := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(btnSelectedAddrCopy, 10, 1, false).
		AddItem(nil, 3, 1, false).
		AddItem(nil, 0, 1, false)

	formDetails := tview.NewForm().
		SetHorizontal(true).
		//AddFormItem(viewSelectedPath).
		//AddFormItem(viewSelectedAddr).
		AddButton("Send", func() {
			if f.selectedAddress != nil {
				f.SwitchToPage(pages.OperationTx, f.selectedAddress)
			}
		}).
		AddButton("Refresh", func() {
			// f.actionUpdateAddresses()
		}).
		AddButton("Paste", func() {
			pasteData, err := clipboard.PasteFromClipboard()
			if err != nil {
				f.Emit(events.EventLogError, fmt.Sprintf("Cannot paste: %s", err))
			} else {
				f.Emit(events.EventLogSuccess, pasteData)
			}
		}).
		AddButton("set W3", func() {
			if f.selectedAddress != nil {
				err := f.API().SetAddressW3(&dto.SetAddressW3DTO{
					MhdaPath: f.selectedAddress.Path,
				})
				if err != nil {
					f.Emit(events.EventLogError, fmt.Sprintf("Cannot set address for Web 3: %s", err))
				} else {
					f.Emit(events.EventLogSuccess, "Address permitted for Web 3")
				}
			}
		})

	formDetails.SetBorderPadding(0, 0, 1, 0)

	formDetails.
		SetTitle(" Send to ").
		SetBorder(true).
		SetBorderColor(tcell.ColorDimGrey)

	formDetails.SetBorderPadding(0, 0, 1, 0)

	formDetails2 := tview.NewForm().
		SetHorizontal(true).
		AddButton("Show QR", func() {
			f.showAddrQR()
		}).
		AddButton("!S node!", func() {
			if f.API() != nil {
				err := f.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					ChainKey:     f.selectedAddress.ChainKey,
					AccountIndex: uint32(f.selectedAddress.Account),
					NodeIndex:    1,
				})
				if err != nil {
					f.Emit(
						events.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					f.Emit(events.EventLog, "SETTED")
				}
			}
		}).
		AddButton("Label set", func() {
			if f.API() != nil {
				err := f.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					ChainKey:     f.selectedAddress.ChainKey,
					AccountIndex: uint32(f.selectedAddress.Account),
					NodeIndex:    1,
				})
				if err != nil {
					f.Emit(
						events.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					f.Emit(events.EventLog, "SETTED")
				}
			}
		}).
		AddButton("!R node!", func() {
			if f.API() != nil {
				err := f.API().RemoveAccountLinkRPC(&dto.RemoveRPCLinkedAccountDTO{
					ChainKey:     f.selectedAddress.ChainKey,
					AccountIndex: uint32(f.selectedAddress.Account),
				})
				if err != nil {
					f.Emit(
						events.EventLogError,
						fmt.Sprintf("Cannot unlink node for account: %s", err),
					)
				} else {
					f.Emit(
						events.EventLog,
						fmt.Sprintf("Unlinked"),
					)
				}
			}
		})

	formDetails2.SetBorderPadding(0, 0, 1, 0)
	formDetails2.
		SetTitle(" Operations ").
		SetBorder(true).
		SetBorderColor(tcell.ColorDimGrey)

	f.layoutAddrSelected.
		AddItem(viewSelectedPath, 1, 1, false).
		AddItem(viewSelectedAddr, 1, 1, false).
		AddItem(layoutSelectedAddrOptions, 0, 1, false)

	f.layout.
		AddItem(f.layoutAddrSelected, 5, 1, false).
		AddItem(formDetails, 5, 1, false).
		AddItem(formDetails2, 3, 1, false).
		AddItem(f.labelQR, 0, 1, false)
	return f.layout
}

func (f *frameAddressesDetailsAddr) showAddrQR() {
	if f.selectedAddress != nil {
		if !f.isQrIsShown {
			// Show
			f.isQrIsShown = true
			f.labelQR.SetTextColor(tcell.ColorBlack).
				SetBackgroundColor(tcell.ColorLightGray)
			f.labelQR.SetText(qrview.NewQrViewText(f.selectedAddress.Address))
		} else {
			// Hide
			f.isQrIsShown = false
			f.labelQR.Clear()
			f.labelQR.SetTextColor(tcell.ColorDefault).
				SetBackgroundColor(tcell.ColorDefault)
		}
	}
}
