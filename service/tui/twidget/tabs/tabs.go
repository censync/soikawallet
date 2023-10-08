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

package tabs

import (
	"fmt"
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type Tabs struct {
	*tview.Flex
	pages    *tview.Pages
	controls *tview.Flex
}

func NewTabs() *Tabs {
	controls := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	pages := tview.NewPages()

	head := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(controls, 0, 1, false).
		AddItem(nil, 0, 2, false)

	tabs := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(head, 3, 1, false).
		AddItem(pages, 0, 1, false)
	return &Tabs{Flex: tabs, controls: controls, pages: pages}
}

func (t *Tabs) AddItem(label string, primitive tview.Primitive) *Tabs {
	name := fmt.Sprintf("tab_%d", t.pages.GetPageCount())
	t.pages.AddPage(name, primitive, true, false)
	btn := tview.NewButton(label).SetSelectedFunc(func() {
		t.pages.SwitchToPage(name)
	})
	btn.SetBackgroundColorActivated(tcell.ColorBlack).
		SetBackgroundColor(tcell.ColorBlack).
		SetLabelColorActivated(tcell.ColorWhite).
		SetStyleAttrs(tcell.AttrBold).
		SetActivatedStyleAttrs(tcell.AttrUnderline | tcell.AttrBold)

	btn.SetBorderColor(tcell.ColorGray).
		SetBorder(true)

	t.controls.AddItem(btn, 0, 1, false)
	if t.pages.GetPageCount() == 1 {
		t.pages.SwitchToPage(label)
	}
	t.pages.SwitchToPage("tab_0")
	return t
}
