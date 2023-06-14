package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type frameDetailsAccount struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedAccount *resp.AccountResponse

	selectedLabelIndex uint32
}

func newFrameDetailsAccount(state *state.State, accountPath *resp.AccountResponse) *frameDetailsAccount {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameDetailsAccount{
		State:           state,
		layout:          layout,
		selectedAccount: accountPath,
	}
}

func (f *frameDetailsAccount) Layout() *tview.Flex {
	label := tview.NewTextView().
		SetText("Account selected")

	label.SetBorderPadding(0, 0, 8, 8)

	accountLabels := f.API().GetAccountLabels()

	inputSelectLabel := tview.NewDropDown().
		SetLabel("Select label")

	for index, title := range accountLabels {
		inputSelectLabel.AddOption(title, func() {
			f.selectedLabelIndex = index
		})
	}

	inputSelectLabel.AddOption(" [ add label ] ", func() {
		f.SwitchToPage(pageNameSettings)
	})

	formAccountDesc := tview.NewForm().
		AddFormItem(inputSelectLabel).
		AddButton("Set label", f.actionSetLabel).
		AddButton("Remove label", f.actionRemoveLabel)

	f.layout.AddItem(label, 1, 1, false).
		AddItem(formAccountDesc, 0, 1, false)
	return f.layout
}

func (f *frameDetailsAccount) actionSetLabel() {
	err := f.API().SetLabelLink(&dto.SetLabelLinkDTO{
		LabelType: types.AccountLabel,
		Index:     f.selectedLabelIndex,
		Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(event_bus.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot set label: %s", err))
	}
}

func (f *frameDetailsAccount) actionRemoveLabel() {
	err := f.API().RemoveLabelLink(&dto.RemoveLabelLinkDTO{
		LabelType: types.AccountLabel,
		Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(event_bus.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot remvove label: %s", err))
	}
}
