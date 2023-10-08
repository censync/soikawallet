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

package twidget

import "github.com/censync/tview"

type BaseFrame struct {
	baseLayout *tview.Flex
}

func NewBaseFrame(layout *tview.Flex) *BaseFrame {
	return &BaseFrame{baseLayout: layout}
}
func (b *BaseFrame) BaseLayout() *tview.Flex { return b.baseLayout }

func (b *BaseFrame) Layout() *tview.Flex { return b.baseLayout }

func (b *BaseFrame) FuncOnShow() {}

func (b *BaseFrame) FuncOnHide() {
	if b.baseLayout != nil {
		b.baseLayout.Clear()
	}
}

func (b *BaseFrame) FuncOnDraw() {}
