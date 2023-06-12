package walletframe

import (
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type frameDetailsAccount struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedCoinType types.CoinType
	selectedAccount  types.AccountIndex
}

func newFrameDetailsAccount(state *state.State, coinType types.CoinType, accountIndex types.AccountIndex) *frameDetailsAccount {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameDetailsAccount{
		State:            state,
		layout:           layout,
		selectedCoinType: coinType,
		selectedAccount:  accountIndex,
	}
}

func (f *frameDetailsAccount) Layout() *tview.Flex {
	label := tview.NewTextView().
		SetText("Account selected")
	label.SetBorderPadding(0, 0, 8, 8)
	f.layout.AddItem(nil, 0, 1, false).
		AddItem(label, 0, 1, false).
		AddItem(nil, 0, 1, false)
	return f.layout
}
