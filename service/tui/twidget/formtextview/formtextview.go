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

package formtextview

import (
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type FormTextView struct {
	*tview.TextView
	label *tview.TextView

	width  int
	height int
}

func NewFormTextView(value string) *FormTextView {
	formTextView := &FormTextView{TextView: tview.NewTextView()}
	formTextView.SetText(value)
	formTextView.SetDynamicColors(true)
	//_, _, _, formTextView.height = formTextView.TextView.GetRect()
	formTextView.height = 1
	return formTextView
}

// Primitive

// FormItem
func (t *FormTextView) GetLabel() string {
	return ``
}

func (t *FormTextView) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	return nil
}

func (t *FormTextView) GetFieldWidth() int {
	return t.width
}

func (t *FormTextView) GetFieldHeight() int {
	return t.height
}

func (t *FormTextView) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	return t
}
