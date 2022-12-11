package walletframe

import (
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
)

type pageAbout struct {
	*BaseFrame
	*state.State
}

func newPageAbout(state *state.State) *pageAbout {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(5, 0, 0, 0)

	return &pageAbout{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageAbout) FuncOnShow() {
	layoutAbout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	viewInfo := tview.NewTextView().
		// SetWordWrap(false).
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter).
		SetText(p.Tr().T("ui.label", "info"))

	viewContacts := tview.NewTextView().
		SetScrollable(false).
		SetText("Website: https://soikawallet.app\n\nGitHub: https://github.com/censync/soikawallet\n\nTwitter: https://twitter.com/SoikaWallet")

	layoutAbout.AddItem(viewInfo, 0, 1, false).
		AddItem(viewContacts, 0, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutAbout, 0, 2, false).
		AddItem(nil, 0, 1, false)
}
