package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type frameDetailsAddr struct {
	layout *tview.Flex
	*state.State

	// ui
	layoutAddrSelected *tview.Flex
	labelQR            *tview.TextView

	// vars
	selectedAddress *resp.AddressResponse
}

func newFrameDetailsAddr(state *state.State, selectedAddress *resp.AddressResponse) *frameDetailsAddr {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	return &frameDetailsAddr{
		State:           state,
		layout:          layout,
		selectedAddress: selectedAddress,
	}
}

func (f *frameDetailsAddr) Layout() *tview.Flex {

	f.layoutAddrSelected = tview.NewFlex().
		SetDirection(tview.FlexRow)

	//SetBorderPadding(0, 0, 0, 0)
	f.layoutAddrSelected.SetTitle(" Send from ").
		SetBorder(true).
		SetBorderColor(tcell.ColorDimGrey).
		SetBorderPadding(0, 0, 3, 0)

	f.labelQR = tview.NewTextView().
		SetWordWrap(false).
		//SetTextColor(tcell.ColorBlack).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	// Selected address label
	// clipboard.CopyToClipboard(inputMnemonic.GetText())
	viewSelectedPath := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorLightBlue).
		SetText(f.selectedAddress.Path)

	viewSelectedAddr := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorBlue).
		SetText(f.selectedAddress.Address)

	btnSelectedAddrCopy := tview.NewButton("copy").SetSelectedFunc(func() {
		if f.selectedAddress != nil {
			err := clipboard.CopyToClipboard(f.selectedAddress.Address)
			if err != nil {
				f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot copy to clipboard: %s", err))
			} else {
				f.Emit(event_bus.EventLogSuccess, "Address copied to clipboard")
			}

		}
	})

	btnSelectedAddrSetW3 := tview.NewButton("Set W3").SetSelectedFunc(func() {
		if f.selectedAddress != nil {
			err := f.API().SetAddressW3(&dto.SetAddressW3DTO{
				DerivationPath: f.selectedAddress.Path,
			})
			if err != nil {
				f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot set address for Web 3: %s", err))
			} else {
				f.Emit(event_bus.EventLogSuccess, "Address permitted for Web 3")
			}
		}
	})

	layoutSelectedAddrOptions := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(btnSelectedAddrCopy, 10, 1, false).
		AddItem(nil, 3, 1, false).
		AddItem(btnSelectedAddrSetW3, 10, 1, false).
		AddItem(nil, 0, 1, false)

	formDetails := tview.NewForm().
		SetHorizontal(true).
		//AddFormItem(viewSelectedPath).
		//AddFormItem(viewSelectedAddr).
		AddButton("Send", func() {
			if f.selectedAddress != nil {
				f.SwitchToPage(pageNameOperationTx, f.selectedAddress)
			}
		}).
		AddButton("Refresh", func() {
			//go f.actionUpdateAddresses()
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
					CoinType:     uint32(f.selectedAddress.CoinType),
					AccountIndex: uint32(f.selectedAddress.Account),
					NodeIndex:    1,
				})
				if err != nil {
					f.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					f.Emit(
						event_bus.EventLog,
						fmt.Sprintf("SETTED"),
					)
				}
			}
		}).
		AddButton("Label set", func() {
			if f.API() != nil {
				err := f.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					CoinType:     uint32(f.selectedAddress.CoinType),
					AccountIndex: uint32(f.selectedAddress.Account),
					NodeIndex:    1,
				})
				if err != nil {
					f.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					f.Emit(
						event_bus.EventLog,
						fmt.Sprintf("SETTED"),
					)
				}
			}
		}).
		AddButton("!R node!", func() {
			if f.API() != nil {
				err := f.API().RemoveAccountLinkRPC(&dto.RemoveRPCLinkedAccountDTO{
					CoinType:     uint32(f.selectedAddress.CoinType),
					AccountIndex: uint32(f.selectedAddress.Account),
				})
				if err != nil {
					f.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot unlink node for account: %s", err),
					)
				} else {
					f.Emit(
						event_bus.EventLog,
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
		//AddItem(tview.NewTextView().SetText("[ Send to ]").SetTextColor(tcell.ColorYellow), 1, 1, false).
		AddItem(formDetails, 5, 1, false).
		AddItem(formDetails2, 3, 1, false).
		AddItem(f.labelQR, 0, 1, false)
	return f.layout
}

func (f *frameDetailsAddr) clearAddrQR() {
	// f.labelQR.Clear()
	// f.labelQR.SetTextColor(tcell.ColorDefault).
	//	SetBackgroundColor(tcell.ColorDefault)
}

func (f *frameDetailsAddr) showAddrQR() {
	//p.Emit(handler.EventShowModal, modalQR)
}
