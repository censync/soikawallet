package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/wallet"
	"github.com/censync/soikawallet/util/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
	"strings"
)

type pageRestoreMnemonic struct {
	*BaseFrame
	*state.State

	// ui
	inputMnemonic *tview.TextArea
	inputPassword *tview.InputField
}

func newPageRestoreMnemonic(state *state.State) *pageRestoreMnemonic {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageRestoreMnemonic{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageRestoreMnemonic) FuncOnShow() {
	p.inputMnemonic = tview.NewTextArea()
	p.inputMnemonic.SetTitle(`Mnemonic`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	p.inputPassword = tview.NewInputField().
		SetMaskCharacter('*')
	p.inputPassword.SetTitle(`Password`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	// env
	envMnemonic, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC")
	if ok {
		p.inputMnemonic.SetText(strings.TrimSpace(envMnemonic), true)
		os.Unsetenv("SOIKAWALLET_MNEMONIC")
	}

	envMnemonicPassphrase, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC_PASSPHRASE")
	if ok {
		p.inputPassword.SetText(strings.TrimSpace(envMnemonicPassphrase))
		os.Unsetenv("SOIKAWALLET_MNEMONIC_PASSPHRASE")
	}

	btnNext := tview.NewButton(p.Tr().T("ui.button", "next")).
		SetStyleAttrs(tcell.AttrBold).
		SetSelectedFunc(p.actionRestoreWithMnemonic)

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

	layoutRestoreForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.inputMnemonic, 10, 1, false).
		AddItem(p.inputPassword, 3, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutRestoreForm, 0, 3, false).
		AddItem(layoutWizard, 35, 1, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageRestoreMnemonic) actionRestoreWithMnemonic() {
	instanceId, err := wallet.API().Init(&dto.InitWalletDTO{
		Mnemonic:   p.inputMnemonic.GetText(),
		Passphrase: p.inputPassword.GetText(),
	})
	if err != nil {
		p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot restore wallet: %s", err))
	} else {
		//p.SetWallet(wallet)
		p.Emit(event_bus.EventUpdateCurrencies, nil)
		clipboard.Clear()
		p.Emit(event_bus.EventWalletInitialized, instanceId)
		p.SwitchToPage(pageNameCreateWallets)
	}
}

func (p *pageRestoreMnemonic) FuncOnHide() {
	p.layout.Clear()
}
