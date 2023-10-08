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

package strip_color

import (
	"github.com/gdamore/tcell/v2"
	"testing"
)

func TestStripColor_Next(t *testing.T) {
	colors := []tcell.Color{tcell.ColorBlack, tcell.ColorGrey, tcell.ColorLightGray}
	s := NewStripColor(colors...)

	for attempts := 3; attempts > 0; attempts-- {
		for i := 0; i < len(colors); i++ {
			if s.index != i {
				t.Fatal("incorrect index")
			}
			s.Next()
		}
	}

}

func TestStripColor_Flush(t *testing.T) {
	s := NewStripColor(tcell.ColorBlack, tcell.ColorGrey)
	s.Next()
	s.Flush()
	if s.index != 0 {
		t.Fatal("incorrect index")
	}
}
