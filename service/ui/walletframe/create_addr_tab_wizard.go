package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

func (p *pageCreateWallet) tabWizard() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	p.layoutCreateWalletsForm = tview.NewFlex().
		SetDirection(tview.FlexRow)

	layoutGlobalSettings := tview.NewForm().
		SetHorizontal(true)

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetOptions(types.GetCoinNames(), func(text string, index int) {
			p.selectedChain = types.GetCoinByName(text)
		}).
		SetCurrentOption(0)

	// inputSelectNetwork.SetBorderPadding(1, 1, 1, 1)

	inputUseHardenedAddresses := tview.NewCheckbox().
		SetLabel("Use hardened").
		SetChangedFunc(func(checked bool) {
			p.selectedUseHardened = checked
		})

	layoutGlobalSettings.
		AddFormItem(inputSelectNetwork).
		AddFormItem(inputUseHardenedAddresses)

	p.layoutCreateWalletsForm.
		AddItem(layoutGlobalSettings, 3, 1, true)

	for addressIndex := 0; addressIndex < addressPoolGap; addressIndex++ {
		labelWalletForm := tview.NewForm().
			SetHorizontal(true).
			SetItemPadding(1).
			AddInputField("Account", "1", 10, tview.InputFieldInteger, nil).
			AddDropDown("Charge", []string{"External", "Internal"}, 0, func(text string, optionIndex int) {
				if optionIndex == 0 {
					p.selectedCharge = 0
				} else {
					p.selectedCharge = 1
				}
			}).
			AddInputField("AddressIndex", fmt.Sprintf("%d", addressIndex), 10, tview.InputFieldInteger, nil)

		p.layoutCreateWalletsForm.AddItem(labelWalletForm, 3, 1, false)
	}

	labelButtons := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "create"), func() {
			p.actionCreateAddrWizard()
		})

	layout.AddItem(p.layoutCreateWalletsForm, 90, 1, false).
		AddItem(labelButtons, 10, 1, false)
	return layout
}

func (p *pageCreateWallet) actionCreateAddrWizard() {
	req := &dto.AddAddressesDTO{}

	// Skip zero index for dropdown and checkbox
	for entry := 1; entry < p.layoutCreateWalletsForm.GetItemCount(); entry++ {
		entryItem := p.layoutCreateWalletsForm.GetItem(entry).(*tview.Form)
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
		p.Emit(handler.EventLogError, fmt.Sprintf("Cannot create addresses: %s", err))
	} else {
		for _, addr := range addresses {
			p.Emit(handler.EventLogInfo, fmt.Sprintf("Added address: %s %s", addr.Path, addr.Address))
		}
		p.SwitchToPage(pageNameAddresses)
	}
}
