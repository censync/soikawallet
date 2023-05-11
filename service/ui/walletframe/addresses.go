package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/spinner"
	"github.com/censync/soikawallet/service/ui/widgets/strip_color"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
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
	addrTree            *tview.TreeNode
	layoutDetails       *tview.Flex
	layoutAddrSelected  *tview.Flex
	labelQR             *tview.TextView

	// var
	selectedAddress *resp.AddressResponse
	isUpdating      bool
	balanceSpinner  *spinner.Spinner
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
	p.layoutDetails = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutAddrSelected = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutAddrSelected.SetTitle(" Send from ").
		SetBorder(true).
		SetBorderColor(tcell.ColorDimGrey).
		SetBorderPadding(0, 0, 3, 0)

	p.labelQR = tview.NewTextView().
		SetWordWrap(false).
		//SetTextColor(tcell.ColorBlack).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter)

	//SetBorderPadding(0, 0, 0, 0)

	// Selected address label
	// clipboard.CopyToClipboard(inputMnemonic.GetText())
	viewSelectedPath := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorLightBlue)
	viewSelectedAddr := tview.NewTextView().
		SetToggleHighlights(true).
		SetTextColor(tcell.ColorBlue)

	btnSelectedAddrCopy := tview.NewButton("copy").SetSelectedFunc(func() {
		if p.selectedAddress != nil {
			err := clipboard.CopyToClipboard(p.selectedAddress.Address)
			if err != nil {
				p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot copy to clipboard: %s", err))
			} else {
				p.Emit(event_bus.EventLogSuccess, "Address copied to clipboard")
			}

		}
	})

	layoutSelectedAddrOptions := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(btnSelectedAddrCopy, 10, 1, false).
		AddItem(nil, 0, 1, false)

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
			} else if action == tview.MouseLeftClick && p.selectedAddress != nil {
				if p.layoutAddressesTree.GetCurrentNode().GetLevel() != addrNodeLevelAddr {
					p.selectedAddress = nil
					p.clearLayoutSelected()
				}
			}
		}
		return action, event
	})

	p.layoutAddressesTree.SetBorder(true)

	p.layoutAddressesTree.SetSelectedFunc(func(node *tview.TreeNode) {
		p.clearLayoutSelected()

		reference := node.GetReference()
		if reference != nil {
			if addressEntry, ok := reference.(*addrNodeViewEntry); ok {
				p.selectedAddress = addressEntry.addr
			}
		}
		if p.selectedAddress != nil {
			p.updateLayoutAddrDetails()
		}

	})

	formDetails := tview.NewForm().
		SetHorizontal(true).
		//AddFormItem(viewSelectedPath).
		//AddFormItem(viewSelectedAddr).
		AddButton("Send", func() {
			if p.selectedAddress != nil {
				p.SwitchToPage(pageNameOperationTx, p.selectedAddress)
			}
		}).
		AddButton("Refresh", func() {
			go p.actionUpdateAddresses()
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
			p.showAddrQR()
		}).
		AddButton("!S node!", func() {
			if p.API() != nil {
				err := p.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					CoinType:     60,
					AccountIndex: 1,
					NodeIndex:    1,
				})
				if err != nil {
					p.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					p.Emit(
						event_bus.EventLog,
						fmt.Sprintf("SETTED"),
					)
				}
			}
		}).
		AddButton("Label set", func() {
			if p.API() != nil {
				err := p.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					CoinType:     60,
					AccountIndex: 1,
					NodeIndex:    1,
				})
				if err != nil {
					p.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					p.Emit(
						event_bus.EventLog,
						fmt.Sprintf("SETTED"),
					)
				}
			}
		}).
		AddButton("!R node!", func() {
			if p.API() != nil {
				err := p.API().RemoveAccountLinkRPC(&dto.RemoveRPCLinkedAccountDTO{
					CoinType:     60,
					AccountIndex: 1,
				})
				if err != nil {
					p.Emit(
						event_bus.EventLogError,
						fmt.Sprintf("Cannot unlink node for account: %s", err),
					)
				} else {
					p.Emit(
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

	p.layoutAddrSelected.
		AddItem(viewSelectedPath, 1, 1, false).
		AddItem(viewSelectedAddr, 1, 1, false).
		AddItem(layoutSelectedAddrOptions, 0, 1, false)

	p.layoutDetails.
		AddItem(p.layoutAddrSelected, 5, 1, false).
		//AddItem(tview.NewTextView().SetText("[ Send to ]").SetTextColor(tcell.ColorYellow), 1, 1, false).
		AddItem(formDetails, 5, 1, false).
		AddItem(formDetails2, 3, 1, false).
		AddItem(p.labelQR, 0, 1, false)

	p.layout.
		AddItem(p.layoutAddressesTree, 0, 1, false).
		AddItem(p.layoutDetails, 50, 1, false)

	return p.layout
}

func (p *pageAddresses) FuncOnShow() {
	go p.actionUpdateAddresses()
}

func (p *pageAddresses) actionUpdateAddresses() {
	if p.isUpdating {
		p.Emit(
			event_bus.EventLog,
			fmt.Sprintf("Updating in process"),
		)
		return
	}

	p.isUpdating = true

	p.selectedAddress = nil
	p.addrTree.ClearChildren()
	p.clearLayoutSelected()

	if p.API() != nil {
		for _, coin := range types.GetCoinTypes() {
			accounts := p.API().GetAccountsByCoin(&dto.GetAccountsByCoinDTO{
				CoinType: uint32(coin),
			})

			if len(accounts) == 0 {
				continue
			}

			coinNode := tview.NewTreeNode(types.GetCoinNameByIndex(coin))

			for _, accountIndex := range accounts {
				accountNode := tview.NewTreeNode(strconv.Itoa(int(accountIndex)))
				stripColor := strip_color.NewStripColor(tcell.ColorLightGray, tcell.ColorDimGrey)
				for _, address := range p.API().GetAddressesByAccount(&dto.GetAddressesByAccountDTO{
					CoinType:     uint32(coin),
					AccountIndex: uint32(accountIndex),
				}) {
					addressNode := tview.NewTreeNode(fmt.Sprintf(
						"%d - %s",
						address.AddressIndex.Index,
						p.addrTruncate(address.Address),
					),
					)
					addressNode.SetReference(&addrNodeViewEntry{
						addr: address,
					})
					addressNode.SetColor(stripColor.Next())
					accountNode.AddChild(addressNode)
				}
				coinNode.AddChild(accountNode)
			}
			p.addrTree.AddChild(coinNode)
		}
		p.Emit(event_bus.EventDrawForce, nil)

		p.balanceSpinner.Start(p.actionTreeSpinnersUpdateFrame)
		p.actionUpdateBalances()
	}
}

func (p *pageAddresses) actionUpdateBalances() {
	for _, coinTree := range p.addrTree.GetChildren() {
		for _, accountTree := range coinTree.GetChildren() {
			for _, addrTree := range accountTree.GetChildren() {
				if addrTree.GetReference() != nil {
					addrView := addrTree.GetReference().(*addrNodeViewEntry)
					if addrView.balances == nil {
						balances, err := p.API().GetTokensBalancesByPath(&dto.GetAddressTokensByPathDTO{
							DerivationPath: addrView.addr.Path,
						})
						//p.Emit(handler.EventLog, "actionUpdateBalances get data")
						if err != nil {
							p.Emit(
								event_bus.EventLogError,
								fmt.Sprintf("Cannot get data for %s: %s", addrView.addr.Path, err),
							)
						} else {
							balancesStr := ""
							for key, value := range balances {
								balancesStr += fmt.Sprintf("$%s - %f ", key, value)
							}
							addrTree.SetText(fmt.Sprintf(
								"%d - %s | %s",
								addrView.addr.AddressIndex.Index,
								p.addrTruncate(addrView.addr.Address), // format long addr
								balancesStr,
							))
							x := 22
							addrView.balances = &x
						}
					}
				}
			}
		}
	}
}

func (p *pageAddresses) actionTreeSpinnersUpdateFrame(frame string) {
	var isSpinnable bool
	for _, coinTree := range p.addrTree.GetChildren() {
		for _, accountTree := range coinTree.GetChildren() {
			for _, addrTree := range accountTree.GetChildren() {
				if addrTree.GetReference() != nil {
					addrView := addrTree.GetReference().(*addrNodeViewEntry)
					if addrView.balances == nil {
						isSpinnable = true
						// TODO: mutex or duplicate view required
						addrTree.SetText(fmt.Sprintf(
							"%d - %s | %s",
							addrView.addr.AddressIndex.Index,
							p.addrTruncate(addrView.addr.Address),
							frame,
						))
					}
				}
			}
		}
	}
	p.Emit(event_bus.EventDrawForce, nil)
	if !isSpinnable {
		p.isUpdating = false
		p.balanceSpinner.Stop()
	}
}

func (p *pageAddresses) clearLayoutSelected() {
	p.clearAddrQR()
	// path
	p.layoutAddrSelected.GetItem(0).(*tview.TextView).Clear()
	// addr
	p.layoutAddrSelected.GetItem(1).(*tview.TextView).Clear()
	// btn copy
	p.layoutAddrSelected.GetItem(2).(*tview.Flex).
		GetItem(0).(*tview.Button).SetDisabled(true)
}

func (p *pageAddresses) updateLayoutAddrDetails() {
	p.clearLayoutSelected()
	p.layoutAddrSelected.GetItem(0).(*tview.TextView).
		SetText(p.selectedAddress.Path)
	p.layoutAddrSelected.GetItem(1).(*tview.TextView).
		SetText(p.selectedAddress.Address)
	p.layoutAddrSelected.GetItem(2).(*tview.Flex).
		GetItem(0).(*tview.Button).SetDisabled(false)
}

func (p *pageAddresses) clearAddrQR() {
	p.labelQR.Clear()
	p.labelQR.SetTextColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)
}

func (p *pageAddresses) showAddrQR() {
	//p.Emit(handler.EventShowModal, modalQR)
}

func (p *pageAddresses) addrTruncate(src string) string {
	x1, _, x2, _ := p.layoutAddressesTree.GetInnerRect()
	if len(src) <= 14 || x2-x1 > 60 {
		return src
	}
	return src[:6] + "..." + src[len(src)-5:]
}
