package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
	"os"
	"strings"
)

type pageRestoreMnemonic struct {
	*BaseFrame
	*state.State
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
	inputMnemonic := tview.NewTextArea()
	inputMnemonic.SetTitle(`Mnemonic`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)
	inputPassword := tview.NewInputField().
		SetMaskCharacter('*')
	inputPassword.SetTitle(`Password`).
		SetTitleAlign(tview.AlignLeft).
		SetBorder(true)

	// env
	envMnemonic, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC")
	if ok {
		inputMnemonic.SetText(strings.TrimSpace(envMnemonic), true)
	}

	envMnemonicPassphrase, ok := os.LookupEnv("SOIKAWALLET_MNEMONIC_PASSPHRASE")
	if ok {
		inputPassword.SetText(strings.TrimSpace(envMnemonicPassphrase))
	}

	labelNext := tview.NewForm().
		SetHorizontal(true).
		SetItemPadding(1).
		AddButton(p.Tr().T("ui.button", "next"), func() {
			err := service.API().Init(&dto.InitWalletDTO{
				Mnemonic:   inputMnemonic.GetText(),
				Passphrase: inputPassword.GetText(),
			})
			if err != nil {
				p.Emit(handler.EventLogError, fmt.Sprintf("Cannot restore wallet: %s", err))
			} else {
				//p.SetWallet(wallet)
				p.SwitchToPage(pageNameCreateWallets)
			}
		})

	layoutRestoreForm := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputMnemonic, 10, 1, false).
		AddItem(inputPassword, 3, 1, false)

	p.layout.AddItem(layoutRestoreForm, 40, 1, false).
		AddItem(labelNext, 20, 1, false)
}

func (p *pageRestoreMnemonic) FuncOnHide() {
	p.layout.Clear()
}
