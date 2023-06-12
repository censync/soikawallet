package walletframe

import (
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
)

type frameDetailsEmpty struct {
	layout *tview.Flex
	*state.State
}

func newFrameDetailsEmpty(state *state.State) *frameDetailsEmpty {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameDetailsEmpty{
		State:  state,
		layout: layout,
	}
}

func (f *frameDetailsEmpty) Layout() *tview.Flex {
	label := tview.NewTextView().
		SetText("Account or address not selected")
	label.SetBorderPadding(0, 0, 8, 8)
	f.layout.AddItem(nil, 0, 1, false).
		AddItem(label, 0, 1, false).
		AddItem(nil, 0, 1, false)
	return f.layout
}
