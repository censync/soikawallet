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

package create_addresses

import (
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

func (p *pageCreateAddr) tabBulk() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	p.inputDerivationPaths = tview.NewTextArea()
	p.inputDerivationPaths.SetTitle(`Derivation addresses list`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	p.inputDerivationPaths.
		SetPlaceholderStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGrey)).
		SetPlaceholder("m/44'/60'/0'/0/0\nm/44'/60'/0'/0/1\nm/44'/60'/0'/0/2")

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(p.actionCreateAddrBulk)

	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnNext, 3, 1, false)

	layoutWizard.SetBorderPadding(1, 1, 2, 2)

	layout.AddItem(p.inputDerivationPaths, 0, 1, false).
		AddItem(layoutWizard, 35, 1, false)
	return layout
}

func (p *pageCreateAddr) actionCreateAddrBulk() {

	// TODO: Update regexp
	/*
		req := &dto.AddAddressesDTO{}
		bulkStr := p.rxAddressPath.FindAllString(p.inputDerivationPaths.GetText(), 1000)
		sort.Strings(bulkStr)
		p.inputDerivationPaths.SetText(strings.Join(bulkStr, "\n"), true)

		req.MhdaPaths = bulkStr

		addresses, err := p.API().AddAddresses(req)
		if err != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Cannot create addresses: %s", err))
		} else {
			for _, addr := range addresses {
				p.Emit(events.EventLogInfo, fmt.Sprintf("Added address: %s %s", addr.Path, addr.Address))
			}
			p.SwitchToPage(pages.PageNameAddresses)
		}*/
}
