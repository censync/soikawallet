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

package flexmenu

import (
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type FlexMenu struct {
	largeButtons bool
	*tview.Flex
	items []*menuItem
}

func NewFlexMenu(largeButtons bool) *FlexMenu {
	menuLayout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	menuLayout.Box.SetDontClear(false)
	menuLayout.SetBorderPadding(1, 1, 1, 1)

	return &FlexMenu{
		largeButtons: largeButtons,
		Flex:         menuLayout,
	}
}

type menuItem struct {
	Label  string
	Key    tcell.Key
	Action func()
}

func (i *menuItem) LabelDecorated() string {
	if i.Key == 0 {
		return i.Label
	} else {
		return "[yellow][" + tcell.KeyNames[i.Key] + "] [white]" + i.Label
	}
}

func (f *FlexMenu) AddMenuItem(label string, key tcell.Key, action func()) *FlexMenu {
	item := &menuItem{
		Label:  label,
		Key:    key,
		Action: action,
	}

	f.items = append(f.items, item)

	btn := tview.NewButton(item.LabelDecorated()).
		SetSelectedFunc(item.Action).
		SetLabelAlign(tview.AlignLeft).
		SetActivatedStyleAttrs(tcell.AttrBold)

	btn.SetBorderPadding(0, 0, 2, 0)

	size := 1
	if f.largeButtons {
		size = 3
	}
	f.Flex.AddItem(btn, size, 1, false)
	f.Flex.AddItem(nil, 1, 1, false)

	return f
}
