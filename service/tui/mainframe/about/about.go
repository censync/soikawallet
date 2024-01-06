// Copyright 2023 The soikawallet Authors
// This file is part of soikawallet.
//
// soikawallet is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// soikawallet is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package about

import (
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageAbout struct {
	*twoframes.BaseFrame
	*state.State
}

func NewPageAbout(state *state.State) *pageAbout {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(5, 0, 0, 0)

	return &pageAbout{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
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
		SetText("Website: https://soikawallet.app\n\nGitHub: https://github.com/censync/soikawallet\n\nTwitter: https://twitter.com/SoikaWallet\n\n\n\nCreated by immigrants with ♥")

	layoutAbout.AddItem(viewInfo, 0, 1, false).
		AddItem(viewContacts, 0, 1, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(layoutAbout, 0, 2, false).
		AddItem(nil, 0, 1, false)
}
