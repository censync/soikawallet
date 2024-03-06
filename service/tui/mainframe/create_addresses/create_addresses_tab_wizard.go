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
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

const defaultAddrPoolGap = 5

func (p *pageCreateAddr) tabWizard() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	p.layoutAddrEntriesForm = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutAddrEntriesForm.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(` ` + p.Tr().T("ui.label", "addresses") + ` `)

	p.actionUpdateForm()

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(p.actionCreateAddrWizard)

	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnNext, 3, 1, false)

	layoutWizard.SetBorderPadding(1, 1, 2, 2)

	layout.
		AddItem(p.uiGlobalSettingsForm(), 40, 1, false).
		AddItem(p.layoutAddrEntriesForm, 70, 1, false).
		AddItem(layoutWizard, 35, 1, false)
	return layout
}

func (p *pageCreateAddr) uiGlobalSettingsForm() *tview.Form {
	layoutGlobalSettings := tview.NewForm().
		SetHorizontal(false)

	layoutGlobalSettings.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(` ` + p.Tr().T("ui.label", "options") + ` `)

	layoutGlobalSettings.SetBorderPadding(0, 1, 3, 1)

	inputSelectNetwork := tview.NewDropDown().
		SetLabel(p.Tr().T("ui.label", "choose_chain")).
		SetFieldWidth(10).
		// TODO: Optimize it
		SetOptions(p.API().GetAllChainNames(), func(text string, index int) {
			p.selectedChain = p.API().GetChainByName(&dto.GetChainByNameDTO{
				ChainName: text,
			})
			p.actionUpdateSelectedChain()
		}).
		SetCurrentOption(0)

	p.inputSelectDerivationType = tview.NewDropDown().
		SetLabel(p.Tr().T("ui.label", "derivation_type")).
		SetFieldWidth(10).
		// []string{"Root", "BIP-32", "BIP-44"}
		SetOptions([]string{"BIP-44"}, func(text string, index int) {
			/*switch index {
			case 0:
				p.selectedDerivationType = mhda.ROOT
			case 1:
				p.selectedDerivationType = mhda.BIP32
			case 2:
				p.selectedDerivationType = mhda.BIP44
			}*/
			p.selectedDerivationType = mhda.BIP44

			p.actionUpdateSelectedChain()
		}).
		SetCurrentOption(0)

	p.inputSelectDerivationPath = tview.NewTextView().
		//SetLabel(p.Tr().T("ui.label", "derivation_path")).
		SetText("m/44'/coin'/account'/charge/address")

	p.inputUseHardenedAddresses = tview.NewCheckbox().
		SetLabel(p.Tr().T("ui.label", "use_hardened")).
		SetChangedFunc(func(checked bool) {
			p.selectedUseHardened = checked // fix checkbox for return to page
		})

	p.inputAccountIndex = tview.NewInputField().
		SetLabel(`Start account`).
		SetFieldWidth(10).
		SetText(strconv.Itoa(p.accountStartIndex)).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetChangedFunc(func(text string) {
			startIndex, err := strconv.Atoi(text)
			if err == nil && startIndex >= 0 {
				p.accountStartIndex = startIndex
				p.actionUpdateForm()
			}
		})

	p.inputAddrIndex = tview.NewInputField().
		SetLabel(`Start addr index`).
		SetFieldWidth(10).
		SetText(strconv.Itoa(p.addrStartIndex)).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetChangedFunc(func(text string) {
			startIndex, err := strconv.Atoi(text)
			if err == nil && startIndex >= 0 {
				p.addrStartIndex = startIndex
				p.actionUpdateForm()
			}
		})

	layoutGlobalSettings.
		AddFormItem(inputSelectNetwork).
		AddFormItem(p.inputSelectDerivationType).
		AddFormItem(p.inputSelectDerivationPath).
		AddFormItem(p.inputUseHardenedAddresses).
		AddFormItem(p.inputAccountIndex).
		AddFormItem(p.inputAddrIndex).
		AddButton(p.Tr().T("ui.label", "row_add"), func() {
			p.addrPoolGap++
			p.actionUpdateForm()
		}).
		AddButton(p.Tr().T("ui.label", "row_remove"), func() {
			if p.addrPoolGap > 1 {
				p.addrPoolGap--
				p.actionUpdateForm()
			}
		})

	return layoutGlobalSettings
}

func (p *pageCreateAddr) actionUpdateForm() {
	p.layoutAddrEntriesForm.Clear()

	maxIndex := p.addrStartIndex + p.addrPoolGap
	for addressIndex := p.addrStartIndex; addressIndex < maxIndex; addressIndex++ {
		labelWalletForm := tview.NewForm().
			SetHorizontal(true).
			SetItemPadding(2).
			AddInputField(p.Tr().T("ui.label", "account"), strconv.Itoa(p.accountStartIndex), 10, tview.InputFieldInteger, nil).
			AddDropDown(p.Tr().T("ui.label", "charge"), []string{" External ▼ ", " Internal "}, 0, func(text string, optionIndex int) {
				if optionIndex == 0 {
					p.selectedCharge = 0
				} else {
					p.selectedCharge = 1
				}
			}).
			AddInputField(p.Tr().T("ui.label", "index"), fmt.Sprintf("%d", addressIndex), 10, tview.InputFieldInteger, nil)
		labelWalletForm.SetBorderPadding(0, 1, 1, 1)
		p.layoutAddrEntriesForm.AddItem(labelWalletForm, 2, 1, false)
	}
}

func (p *pageCreateAddr) actionCreateAddrWizard() {
	req := &dto.AddAddressesDTO{}

	for entry := 0; entry < p.layoutAddrEntriesForm.GetItemCount(); entry++ {
		entryItem := p.layoutAddrEntriesForm.GetItem(entry).(*tview.Form)

		accountStr := entryItem.GetFormItem(0).(*tview.InputField).GetText()
		accountIndex, err := strconv.ParseUint(accountStr, 0, 32)
		if err != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Incorrect account value for row %d", entry))
			return
		}

		indexStr := entryItem.GetFormItem(2).(*tview.InputField).GetText()
		addrIndex, err := strconv.ParseUint(indexStr, 0, 32)
		if err != nil {
			p.Emit(events.EventLogError, fmt.Sprintf("Incorrect index value for row %d", addrIndex))
			return
		}

		dPath := mhda.NewDerivationPath(
			p.selectedDerivationType,
			p.selectedChain.CoinType(),
			mhda.AccountIndex(accountIndex),
			mhda.ChargeType(p.selectedCharge),
			mhda.AddressIndex{
				Index:      uint32(addrIndex),
				IsHardened: p.selectedUseHardened,
			},
		)
		mPath := mhda.NewAddress(p.selectedChain, dPath)
		req.MhdaPaths = append(req.MhdaPaths, mPath.NSS())
	}

	addresses, err := p.API().AddAddresses(req)
	if err != nil {
		p.Emit(events.EventLogError, fmt.Sprintf("Cannot create addresses: %s", err))
	} else {
		for _, addr := range addresses {
			p.Emit(events.EventLogInfo, fmt.Sprintf(
				"Added address: %s %s",
				p.API().GetChainNameByKey(&dto.GetChainNameByKeyDTO{
					ChainKey: p.selectedChain.Key(), //TODO: Optimize it
				}),
				addr.Address,
			),
			)
		}
		p.SwitchToPage(pages.Addresses)
	}
}

func (p *pageCreateAddr) actionUpdateSelectedChain() {
	if p.inputSelectDerivationPath == nil {
		// temp
		return
	}
	switch p.selectedDerivationType {
	case mhda.ROOT:
		// cleanup
		p.inputSelectDerivationPath.SetText("root")
	case mhda.BIP32:
		p.selectedDerivationPath = *mhda.NewDerivationPath(
			mhda.BIP32,
			p.selectedChain.CoinType(),
			mhda.AccountIndex(p.accountStartIndex),
			mhda.ChargeType(p.selectedCharge),
			mhda.AddressIndex{
				Index:      uint32(p.addrStartIndex),
				IsHardened: p.selectedUseHardened,
			})
		p.inputSelectDerivationPath.SetText(p.selectedDerivationPath.String())
	case mhda.BIP44:
		p.selectedDerivationPath = *mhda.NewDerivationPath(
			mhda.BIP44,
			p.selectedChain.CoinType(),
			mhda.AccountIndex(p.accountStartIndex),
			mhda.ChargeType(p.selectedCharge),
			mhda.AddressIndex{
				Index:      uint32(p.addrStartIndex),
				IsHardened: p.selectedUseHardened,
			})
		p.inputSelectDerivationPath.SetText(p.selectedDerivationPath.String())
	}
}
