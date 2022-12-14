package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/spinner"
	"github.com/censync/soikawallet/service/ui/widgets/strip_color"
	"github.com/censync/soikawallet/types"
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
	addrTree       *tview.TreeNode
	layoutDetails  *tview.Flex
	layoutSelected *tview.Flex
	labelQR        *tview.TextView

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
		State:          state,
		BaseFrame:      &BaseFrame{layout: layout},
		balanceSpinner: spinner.NewSpinner(spinner.SpinThree, 180),
	}
}

func (p *pageAddresses) Layout() *tview.Flex {
	p.layoutDetails = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutSelected = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutSelected.SetTitle(" Send from ").
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

	p.addrTree = tview.NewTreeNode("wallets")
	layoutAddressesTree := tview.NewTreeView().
		SetRoot(p.addrTree).SetTopLevel(1)

	// double click for address operations
	layoutAddressesTree.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if layoutAddressesTree.InRect(event.Position()) {
			if action == tview.MouseLeftDoubleClick && p.selectedAddress != nil {
				if layoutAddressesTree.GetCurrentNode().GetLevel() == addrNodeLevelAddr {
					p.SwitchToPage(pageNameOperationTx, p.selectedAddress)
				}
				return action, nil
			}
		}
		return action, event
	})

	layoutAddressesTree.SetBorder(true)

	layoutAddressesTree.SetSelectedFunc(func(node *tview.TreeNode) {
		p.clearLayoutSelected()

		reference := node.GetReference()
		if reference != nil {
			p.selectedAddress = reference.(*resp.AddressResponse)
		}
		if p.selectedAddress != nil {
			p.updateLayoutDetails()
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
		AddButton("!Set node!", func() {
			if p.API() != nil {
				err := p.API().AccountLinkRPCSet(&dto.SetRPCLinkedAccountDTO{
					CoinType:     60,
					AccountIndex: 1,
					NodeIndex:    1,
				})
				if err != nil {
					p.Emit(
						handler.EventLogError,
						fmt.Sprintf("Cannot link node for account: %s", err),
					)
				} else {
					p.Emit(
						handler.EventLog,
						fmt.Sprintf("SETTED"),
					)
				}
			}
		}).
		AddButton("!Remove node!", func() {
			if p.API() != nil {
				err := p.API().RemoveAccountLinkRPC(&dto.RemoveRPCLinkedAccountDTO{
					CoinType:     60,
					AccountIndex: 1,
				})
				if err != nil {
					p.Emit(
						handler.EventLogError,
						fmt.Sprintf("Cannot unlink node for account: %s", err),
					)
				} else {
					p.Emit(
						handler.EventLog,
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

	p.layoutSelected.
		AddItem(viewSelectedPath, 1, 1, false).
		AddItem(viewSelectedAddr, 1, 1, false)

	p.layoutDetails.
		AddItem(p.layoutSelected, 4, 1, false).
		//AddItem(tview.NewTextView().SetText("[ Send to ]").SetTextColor(tcell.ColorYellow), 1, 1, false).
		AddItem(formDetails, 5, 1, false).
		AddItem(formDetails2, 3, 1, false).
		AddItem(p.labelQR, 0, 1, false)

	p.layout.
		AddItem(layoutAddressesTree, 0, 1, false).
		AddItem(p.layoutDetails, 55, 1, false)

	return p.layout
}

func (p *pageAddresses) FuncOnShow() {
	go p.actionUpdateAddresses()
}

func (p *pageAddresses) actionUpdateAddresses() {
	p.selectedAddress = nil
	p.addrTree.ClearChildren()
	p.clearLayoutSelected()
	if p.isUpdating {
		p.Emit(
			handler.EventLog,
			fmt.Sprintf("Update in process"),
		)
		return
	}
	p.isUpdating = true
	defer func() {
		p.Emit(
			handler.EventLog,
			fmt.Sprintf("Update finished"),
		)
		p.isUpdating = false
		p.balanceSpinner.Start(p.actionTreeSpinnersUpdateFrame)
		p.actionUpdateBalances()
	}()

	p.Emit(
		handler.EventLog,
		fmt.Sprintf("Update started"),
	)

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
					/* balancesStr := ""
					balances, err := p.API().GetAddressTokensByPath(&dto.GetAddressTokensByPathDTO{
						DerivationPath: address.Path,
					})

					if err != nil {
						p.Emit(
							handler.EventLogError,
							fmt.Sprintf("Cannot get data for %s: %s", address.Path, err),
						)
					} else {
						for key, value := range balances {
							balancesStr += fmt.Sprintf("$%s - %f ", key, value)
						}
					}

					addressNode := exttree.NewTreeNode(fmt.Sprintf("%d - %s | %s", address.AddressIndex.Index, address.Address, balancesStr))*/
					addressNode := tview.NewTreeNode(fmt.Sprintf("%d - %s", address.AddressIndex.Index, address.Address))
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
		p.Emit(handler.EventDrawForce, nil)
	}
}

func (p *pageAddresses) actionUpdateBalances() {
	for _, coinTree := range p.addrTree.GetChildren() {
		for _, accountTree := range coinTree.GetChildren() {
			for _, addrTree := range accountTree.GetChildren() {
				if addrTree.GetReference() != nil {
					addrView := addrTree.GetReference().(*addrNodeViewEntry)
					if addrView.balances == nil {
						balances, err := p.API().GetAddressTokensByPath(&dto.GetAddressTokensByPathDTO{
							DerivationPath: addrView.addr.Path,
						})
						//p.Emit(handler.EventLog, "actionUpdateBalances get data")
						if err != nil {
							p.Emit(
								handler.EventLogError,
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
								addrView.addr.Address,
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
							addrView.addr.Address,
							frame,
						))
					}
				}
			}
		}
	}
	p.Emit(handler.EventDrawForce, nil)
	if !isSpinnable {
		p.balanceSpinner.Stop()
	}
}

func (p *pageAddresses) clearLayoutSelected() {
	p.clearAddrQR()
	p.layoutSelected.GetItem(0).(*tview.TextView).Clear()
	p.layoutSelected.GetItem(1).(*tview.TextView).Clear()
}

func (p *pageAddresses) updateLayoutDetails() {
	p.clearLayoutSelected()
	p.layoutSelected.GetItem(0).(*tview.TextView).
		SetText(p.selectedAddress.Path)
	p.layoutSelected.GetItem(1).(*tview.TextView).
		SetText(p.selectedAddress.Address)
}

func (p *pageAddresses) clearAddrQR() {
	p.labelQR.Clear()
	p.labelQR.SetTextColor(tcell.ColorDefault).
		SetBackgroundColor(tcell.ColorDefault)
}

func (p *pageAddresses) showAddrQR() {
	//p.Emit(handler.EventShowModal, modalQR)
}
