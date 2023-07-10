package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/wallet"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/censync/soikawallet/util/seed"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

type pageInitMnemonic struct {
	*BaseFrame
	*state.State

	// tui
	inputMnemonic *tview.Form

	// vars
	selectedMnemonicEntropy  int
	selectedMnemonicLanguage string
	mnemonic                 string
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
	p.inputMnemonic = tview.NewForm().
		SetHorizontal(true)

	p.inputMnemonic.SetItemPadding(4).
		SetBorder(true).
		SetTitle(" Mnemonic ").
		SetTitleAlign(tview.AlignLeft)

	inputPassword := tview.NewInputField().
		SetMaskCharacter('*')
	inputPassword.SetTitle(` Password `).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	inputDropDownEntropy := tview.NewDropDown().
		SetLabel("Entropy").
		SetFieldWidth(5).
		SetOptions(seed.EntropyList(), func(option string, optionIndex int) {
			p.selectedMnemonicEntropy, _ = strconv.Atoi(option)
			p.actionMnemonicUpdate()
		}).
		SetCurrentOption(len(seed.EntropyList()) - 1)

	inputDropDownLanguage := tview.NewDropDown().
		SetLabel("Language").
		SetFieldWidth(10).
		SetOptions([]string{"english"}, func(option string, optionIndex int) {
			if p.selectedMnemonicLanguage != option {
				p.selectedMnemonicLanguage = option
				p.actionMnemonicUpdate()
			}
		}).
		SetCurrentOption(0)

	formMnemonicConfig := tview.NewForm().
		SetHorizontal(true).
		AddFormItem(inputDropDownEntropy).
		AddFormItem(inputDropDownLanguage)

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(func() {
			instanceId, err := wallet.API().Init(&dto.InitWalletDTO{
				Mnemonic:   p.mnemonic,
				Passphrase: inputPassword.GetText(),
			})
			if err != nil {
				p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot init wallet: %s", err))
			} else {
				clipboard.Clear()
				p.Emit(event_bus.EventWalletInitialized, instanceId)
				p.SwitchToPage(pageNameCreateWallets)
			}
		})
	btnBack := tview.NewButton(p.Tr().T("ui.button", "back")).
		SetSelectedFunc(func() {
			p.SwitchToPage(p.Pages().GetPrevious())
		})

	layoutWizard := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(btnNext, 3, 1, false).
		AddItem(nil, 1, 1, false).
		AddItem(btnBack, 3, 1, false)

	layoutWizard.SetBorderPadding(1, 1, 2, 2)

	layoutOptions := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "generate_mnemonic"), func() {
			p.actionMnemonicUpdate()
		}).
		AddButton(p.Tr().T("ui.button", "copy_to_clipboard"), func() {
			err := clipboard.CopyToClipboard(p.mnemonic)
			if err != nil {
				p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot copy to clipboard: %s", err))
			}
		})

	layoutMnemonicForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.inputMnemonic, 0, 6, false).
		AddItem(inputPassword, 3, 1, false).
		AddItem(formMnemonicConfig, 3, 1, false).
		AddItem(layoutOptions, 3, 1, false)

	p.layout.AddItem(layoutMnemonicForm, 0, 3, false).
		AddItem(layoutWizard, 35, 1, false)

}

func (p *pageInitMnemonic) actionMnemonicUpdate() {
	var err error
	//p.inputMnemonic.SetText(``, false)
	p.inputMnemonic.Clear(false)

	p.mnemonic, err = wallet.API().GenerateMnemonic(&dto.GenerateMnemonicDTO{
		BitSize:  p.selectedMnemonicEntropy,
		Language: p.selectedMnemonicLanguage,
	})

	if err != nil {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot generate mnemonic: %s", err))
		return
	}

	//p.inputMnemonic.SetText(mnemonic, false)

	mnemonicList := strings.Split(p.mnemonic, ` `)

	for index := range mnemonicList {
		mnemonicInput := tview.NewInputField().
			SetFieldWidth(15).
			SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
				return false
			}).
			SetLabel(fmt.Sprintf("%2d.", index+1)).
			SetText(mnemonicList[index])

		p.inputMnemonic.AddFormItem(mnemonicInput)
	}
}

func (p *pageInitMnemonic) FuncOnHide() {
	p.mnemonic = ``
	p.layout.Clear()
}
