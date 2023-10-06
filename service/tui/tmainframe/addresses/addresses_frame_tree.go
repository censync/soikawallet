package addresses

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/twidget/strip_color"
	"github.com/censync/soikawallet/types"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type accountNodeViewEntry struct {
	account *resp.AccountResponse
}

type addrNodeViewEntry struct {
	addr     *resp.AddressResponse
	balances map[string]float64 //*resp.AddressTokensBalanceListResponse
}

func (p *pageAddresses) actionUpdateAddresses() {
	if p.isUpdating {
		p.Emit(
			events.EventLog,
			fmt.Sprintf("Updating in process"),
		)
		return
	}

	p.isUpdating = true

	p.selectedAddress = nil
	p.addrTree.ClearChildren()

	p.actionUpdateFrameDetails()

	if p.API() != nil {
		for _, chainKey := range types.GetChains() {
			accounts := p.API().GetAccountsByNetwork(&dto.GetAccountsByNetworkDTO{
				ChainKey: chainKey,
			})

			if len(accounts) == 0 {
				continue
			}

			networkNode := tview.NewTreeNode(types.GetNetworkNameByKey(chainKey))

			for _, account := range accounts {
				accountNodeTitle := ""

				if account.Label != "" {
					accountNodeTitle = fmt.Sprintf("[%d] - [blue::b]%s", account.Account, account.Label)
				} else {
					accountNodeTitle = fmt.Sprintf("[%d]", account.Account)
				}

				accountNode := tview.NewTreeNode(accountNodeTitle)

				accountNode.SetReference(&accountNodeViewEntry{
					account: account,
				})
				stripColor := strip_color.NewStripColor(tcell.ColorLightGray, tcell.ColorDimGrey)
				for _, address := range p.API().GetAddressesByAccount(&dto.GetAddressesByAccountDTO{
					ChainKey:     chainKey,
					AccountIndex: uint32(account.Account),
				}) {
					addrIndexFormat := "%d - %s"

					if address.AddressIndex.IsHardened {
						addrIndexFormat = "%d' - %s"
					}

					addressNode := tview.NewTreeNode(fmt.Sprintf(
						addrIndexFormat,
						address.AddressIndex.Index,
						p.addrTruncate(address.Address),
					))
					addressNode.SetReference(&addrNodeViewEntry{
						addr: address,
					})
					if address.IsW3 {
						addressNode.SetColor(tcell.ColorDarkOrange)
					} else {
						addressNode.SetColor(stripColor.Next())
					}

					accountNode.AddChild(addressNode)
				}

				networkNode.AddChild(accountNode)
			}
			p.addrTree.AddChild(networkNode)
		}
		p.Emit(events.EventDrawForce, nil)

		p.balanceSpinner.Start(p.actionTreeSpinnersUpdate)
		p.actionUpdateBalances()
	}
}

func (p *pageAddresses) actionUpdateBalances() {
	for _, networkTree := range p.addrTree.GetChildren() {
		for _, accountTree := range networkTree.GetChildren() {
			for _, addrTree := range accountTree.GetChildren() {
				if addrTree.GetReference() != nil {
					addrView := addrTree.GetReference().(*addrNodeViewEntry)
					if addrView.balances == nil {
						balances, err := p.API().GetTokensBalancesByAddress(&dto.GetAddressTokensByPathDTO{
							MhdaPath: addrView.addr.Path,
						})

						balancesStr := ""
						addrIndexFormat := "%d - %s | %s"
						if err == nil {

							for key, value := range balances {
								balancesStr += fmt.Sprintf("$%s - %f ", key, value)
							}

							if addrView.addr.AddressIndex.IsHardened {
								addrIndexFormat = "%d' - %s | %s"
							}

						} else {
							balances = map[string]float64{} // Empty map for stop anim
							balancesStr = "[gray][cannot get balance]"
							p.Emit(
								events.EventLogError,
								fmt.Sprintf("Cannot get data for %s: %s", addrView.addr.Address, err),
							)
						}
						addrTree.SetText(fmt.Sprintf(
							addrIndexFormat,
							addrView.addr.AddressIndex.Index,
							p.addrTruncate(addrView.addr.Address), // format long addr
							balancesStr,
						))
						addrView.balances = balances
					}
				}
			}
		}
	}
}

func (p *pageAddresses) actionTreeSpinnersUpdate(frame string) {
	var isSpinnable bool
	for _, networkTree := range p.addrTree.GetChildren() {
		for _, accountTree := range networkTree.GetChildren() {
			for _, addrTree := range accountTree.GetChildren() {
				if addrTree.GetReference() != nil {
					addrView := addrTree.GetReference().(*addrNodeViewEntry)
					if addrView.balances == nil {
						isSpinnable = true
						// TODO: mutex or duplicate view required

						addrIndexFormat := "%d - %s | %s"

						if addrView.addr.AddressIndex.IsHardened {
							addrIndexFormat = "%d' - %s | %s"
						}

						addrTree.SetText(fmt.Sprintf(
							addrIndexFormat,
							addrView.addr.AddressIndex.Index,
							p.addrTruncate(addrView.addr.Address),
							frame,
						))
					}
				}
			}
		}
	}
	p.Emit(events.EventDrawForce, nil)
	if !isSpinnable {
		p.isUpdating = false
		p.balanceSpinner.Stop()
	}
}

func (p *pageAddresses) addrTruncate(src string) string {
	x1, _, x2, _ := p.layoutAddressesTree.GetInnerRect()
	if len(src) <= 14 || x2-x1 > 60 {
		return src
	}
	return src[:6] + "..." + src[len(src)-5:]
}
