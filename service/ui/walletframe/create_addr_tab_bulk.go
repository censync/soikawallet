package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

func (p *pageCreateWallet) tabBulk() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	p.inputDerivationPaths = tview.NewTextArea()
	p.inputDerivationPaths.SetTitle(`Derivation addresses list`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	p.inputDerivationPaths.
		SetPlaceholder("m/44'/60'/0'/0/0\nm/44'/60'/0'/0/1\nm/44'/60'/0'/0/2")

	labelButtons := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "create"), func() {
			p.actionCreateAddrBulk()
		})

	layout.AddItem(p.inputDerivationPaths, 0, 1, false).
		AddItem(labelButtons, 20, 1, false)
	return layout
}

func (p *pageCreateWallet) actionCreateAddrBulk() {
	req := &dto.AddAddressesDTO{}

	bulkStr := p.rxAddressPath.FindAllString(p.inputDerivationPaths.GetText(), 1000)
	sort.Strings(bulkStr)
	p.inputDerivationPaths.SetText(strings.Join(bulkStr, "\n"), true)

	req.DerivationPaths = bulkStr

	addresses, err := p.API().AddAddresses(req)
	if err != nil {
		p.Emit(handler.EventLogError, fmt.Sprintf("Cannot create addresses: %s", err))
	} else {
		for _, addr := range addresses {
			p.Emit(handler.EventLogInfo, fmt.Sprintf("Added address: %s %s", addr.Path, addr.Address))
		}
		p.SwitchToPage(pageNameAddresses)
	}
}
