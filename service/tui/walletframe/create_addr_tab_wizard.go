package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
	"strconv"
)

const defaultAddrPoolGap = 5

func (p *pageCreateWallet) tabWizard() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	p.layoutAddrEntriesForm = tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.actionUpdateForm()

	labelButtons := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("tui.button", "create"), func() {
			p.actionCreateAddrWizard()
		})

	layout.
		AddItem(p.uiGlobalSettingsForm(), 40, 1, false).
		AddItem(p.layoutAddrEntriesForm, 70, 1, false).
		AddItem(labelButtons, 20, 1, false)
	return layout
}

func (p *pageCreateWallet) uiGlobalSettingsForm() *tview.Form {
	layoutGlobalSettings := tview.NewForm().
		SetHorizontal(false)

	layoutGlobalSettings.SetBorderPadding(0, 1, 3, 1)

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetFieldWidth(10).
		SetOptions(types.GetCoinNames(), func(text string, index int) {
			p.selectedChain = types.GetCoinByName(text)
		}).
		SetCurrentOption(0)

	inputUseHardenedAddresses := tview.NewCheckbox().
		SetLabel("Use hardened").
		SetChangedFunc(func(checked bool) {
			p.selectedUseHardened = checked
		})

	inputAccountIndex := tview.NewInputField().
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

	inputAddrIndex := tview.NewInputField().
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
		AddFormItem(inputUseHardenedAddresses).
		AddFormItem(inputAccountIndex).
		AddFormItem(inputAddrIndex).
		AddButton("Add row", func() {
			p.addrPoolGap++
			p.actionUpdateForm()
		}).
		AddButton("Remove row", func() {
			if p.addrPoolGap > 1 {
				p.addrPoolGap--
				p.actionUpdateForm()
			}
		})

	return layoutGlobalSettings
}

func (p *pageCreateWallet) actionUpdateForm() {
	p.layoutAddrEntriesForm.Clear()

	maxIndex := p.addrStartIndex + p.addrPoolGap
	for addressIndex := p.addrStartIndex; addressIndex < maxIndex; addressIndex++ {
		labelWalletForm := tview.NewForm().
			SetHorizontal(true).
			SetItemPadding(2).
			AddInputField("Account", strconv.Itoa(p.accountStartIndex), 10, tview.InputFieldInteger, nil).
			AddDropDown("Charge", []string{"External", "Internal"}, 0, func(text string, optionIndex int) {
				if optionIndex == 0 {
					p.selectedCharge = 0
				} else {
					p.selectedCharge = 1
				}
			}).
			AddInputField("Address Index", fmt.Sprintf("%d", addressIndex), 10, tview.InputFieldInteger, nil)
		labelWalletForm.SetBorderPadding(0, 1, 1, 1)
		p.layoutAddrEntriesForm.AddItem(labelWalletForm, 2, 1, false)
	}
}

func (p *pageCreateWallet) actionCreateAddrWizard() {
	req := &dto.AddAddressesDTO{}

	for entry := 0; entry < p.layoutAddrEntriesForm.GetItemCount(); entry++ {
		entryItem := p.layoutAddrEntriesForm.GetItem(entry).(*tview.Form)
		pathFormat := "m/44'/%d'/%s'/%d/%s"
		if p.selectedUseHardened {
			pathFormat += `'`
		}
		dPath := fmt.Sprintf(pathFormat,
			p.selectedChain,
			entryItem.GetFormItem(0).(*tview.InputField).GetText(),
			p.selectedCharge,
			entryItem.GetFormItem(2).(*tview.InputField).GetText(),
		)
		req.DerivationPaths = append(req.DerivationPaths, dPath)
	}

	addresses, err := p.API().AddAddresses(req)
	if err != nil {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot create addresses: %s", err))
	} else {
		for _, addr := range addresses {
			p.Emit(event_bus.EventLogInfo, fmt.Sprintf("Added address: %s %s", addr.Path, addr.Address))
		}
		p.SwitchToPage(pageNameAddresses)
	}
}
