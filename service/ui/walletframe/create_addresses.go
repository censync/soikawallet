package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

const addressPoolGap = 5

type pageCreateWallet struct {
	*BaseFrame
	*state.State
}

func newPageCreateWallet(state *state.State) *pageCreateWallet {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageCreateWallet{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageCreateWallet) FuncOnShow() {
	var (
		selectedChain  types.CoinType
		selectedCharge uint8
	)

	layoutCreateWalletsForm := tview.NewFlex().
		SetDirection(tview.FlexRow)

	layoutGlobalSettings := tview.NewForm().
		SetHorizontal(true)

	inputSelectNetwork := tview.NewDropDown().
		SetLabel("Select network").
		SetOptions(types.GetCoinNames(), func(text string, index int) {
			selectedChain = types.GetCoinByName(text)
		}).
		SetCurrentOption(0)

	// inputSelectNetwork.SetBorderPadding(1, 1, 1, 1)

	inputUseHardenedAddresses := tview.NewCheckbox()
	inputUseHardenedAddresses.SetLabel("Use hardened")

	layoutGlobalSettings.
		AddFormItem(inputSelectNetwork).
		AddFormItem(inputUseHardenedAddresses)

	layoutCreateWalletsForm.
		AddItem(layoutGlobalSettings, 3, 1, true)

	for addressIndex := 0; addressIndex < addressPoolGap; addressIndex++ {
		labelWalletForm := tview.NewForm().
			SetHorizontal(true).
			SetItemPadding(1).
			AddInputField("Account", "1", 10, tview.InputFieldInteger, nil).
			AddDropDown("Charge", []string{"External", "Internal"}, 0, func(text string, optionIndex int) {
				if optionIndex == 0 {
					selectedCharge = 0
				} else {
					selectedCharge = 1
				}
			}).
			AddInputField("AddressIndex", fmt.Sprintf("%d", addressIndex), 10, tview.InputFieldInteger, nil)

		layoutCreateWalletsForm.AddItem(labelWalletForm, 3, 1, false)
	}

	labelButtons := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "create"), func() {
			req := &dto.AddAddressesDTO{}

			// Skip zero index for dropdown and checkbox
			for entry := 1; entry < layoutCreateWalletsForm.GetItemCount(); entry++ {
				entryItem := layoutCreateWalletsForm.GetItem(entry).(*tview.Form)
				pathFormat := "m/44'/%d'/%s'/%d/%s"
				if inputUseHardenedAddresses.IsChecked() {
					pathFormat += `'`
				}
				dPath := fmt.Sprintf(pathFormat,
					selectedChain,
					entryItem.GetFormItem(0).(*tview.InputField).GetText(),
					selectedCharge,
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
		})

	p.layout.AddItem(layoutCreateWalletsForm, 90, 1, false).
		AddItem(labelButtons, 10, 1, false)
}

func (p *pageCreateWallet) FuncOnHide() {
	p.layout.Clear()
}
