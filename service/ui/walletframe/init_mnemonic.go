package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/censync/soikawallet/util/seed"
	"github.com/rivo/tview"
	"strconv"
)

type pageInitMnemonic struct {
	*BaseFrame
	*state.State

	// ui
	inputMnemonic *tview.TextArea

	// vars
	selectedMnemonicEntropy  int
	selectedMnemonicLanguage string
}

func newPageInitMnemonic(state *state.State) *pageInitMnemonic {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageInitMnemonic{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageInitMnemonic) FuncOnShow() {
	p.inputMnemonic = tview.NewTextArea()

	p.inputMnemonic.SetBorder(true).
		SetTitle("Please, save the mnemonic")

	inputPassword := tview.NewInputField().
		SetMaskCharacter('*')
	inputPassword.SetTitle(`Password`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	formMnemonicConfig := tview.NewForm().
		SetHorizontal(true).
		AddDropDown("Entropy", seed.EntropyList(), 0, func(option string, optionIndex int) {
			p.selectedMnemonicEntropy, _ = strconv.Atoi(option)
			p.updateMnemonicConfig()
		}).
		AddDropDown("Language", []string{"english"}, 0, func(option string, optionIndex int) {
			p.selectedMnemonicLanguage = option
			p.updateMnemonicConfig()
		})

	labelNext := tview.NewForm().
		SetHorizontal(false).
		SetItemPadding(2).
		AddButton(p.Tr().T("ui.button", "back"), func() {
			p.SwitchToPage(p.Pages().GetPrevious())
		}).
		AddButton(p.Tr().T("ui.button", "next"), func() {
			err := service.API().Init(&dto.InitWalletDTO{
				Mnemonic:   p.inputMnemonic.GetText(),
				Passphrase: inputPassword.GetText(),
			})
			if err != nil {
				p.Emit(handler.EventLogError, fmt.Sprintf("Cannot init walletInstance: %s", err))
			} else {
				//p.SetWallet(walletInstance)
				p.SwitchToPage(pageNameCreateWallets)
			}
		})

	labelButtons := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "generate_mnemonic"), func() {
			p.updateMnemonicConfig()
		}).
		AddButton(p.Tr().T("ui.button", "copy_to_clipboard"), func() {
			err := clipboard.CopyToClipboard(p.inputMnemonic.GetText())
			if err != nil {
				p.Emit(handler.EventLogError, fmt.Sprintf("Cannot generate mnemonic: %s", err))
			}
		})

	layoutMnemonicForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.inputMnemonic, 10, 1, false).
		AddItem(inputPassword, 3, 1, false).
		AddItem(formMnemonicConfig, 3, 1, false).
		AddItem(labelButtons, 3, 1, false)

	p.layout.AddItem(layoutMnemonicForm, 40, 1, false).
		AddItem(labelNext, 20, 1, false)

}

func (p *pageInitMnemonic) FuncOnHide() {
	p.layout.Clear()
}

func (p *pageInitMnemonic) updateMnemonicConfig() {
	p.inputMnemonic.SetText(``, false)

	mnemonic, err := service.API().GenerateMnemonic(&dto.GenerateMnemonicDTO{
		BitSize:  p.selectedMnemonicEntropy,
		Language: p.selectedMnemonicLanguage,
	})

	if err != nil {
		p.Emit(handler.EventLogError, fmt.Sprintf("Cannot generate mnemonic: %s", err))
	} else {
		p.inputMnemonic.SetText(mnemonic, false)
	}
}
