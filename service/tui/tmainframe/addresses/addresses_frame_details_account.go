package addresses

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/page"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/types"
	"github.com/censync/tview"
)

type frameAddressesDetailsAccount struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedAccount *resp.AccountResponse

	selectedLabelIndex uint32
}

func newAddressesFrameDetailsAccount(state *state.State, accountPath *resp.AccountResponse) *frameAddressesDetailsAccount {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameAddressesDetailsAccount{
		State:           state,
		layout:          layout,
		selectedAccount: accountPath,
	}
}

func (f *frameAddressesDetailsAccount) Layout() *tview.Flex {
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
		f.SwitchToPage(page.Settings)
	})

	formAccountDesc := tview.NewForm().
		AddFormItem(inputSelectLabel).
		AddButton("set label", f.actionSetLabel).
		AddButton("Remove label", f.actionRemoveLabel)

	f.layout.AddItem(label, 1, 1, false).
		AddItem(formAccountDesc, 0, 1, false)
	return f.layout
}

func (f *frameAddressesDetailsAccount) actionSetLabel() {
	err := f.API().SetLabelLink(&dto.SetLabelLinkDTO{
		LabelType: types.AccountLabel,
		Index:     f.selectedLabelIndex,
		// TODO: Finish
		//Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(event_bus.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot set label: %s", err))
	}
}

func (f *frameAddressesDetailsAccount) actionRemoveLabel() {
	err := f.API().RemoveLabelLink(&dto.RemoveLabelLinkDTO{
		LabelType: types.AccountLabel,
		// TODO: Finish
		// Path:      f.selectedAccount.Path,
	})
	if err == nil {
		f.Emit(event_bus.EventLogSuccess, "Label saved for account")
	} else {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot remvove label: %s", err))
	}
}
