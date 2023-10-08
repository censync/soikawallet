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

import "github.com/gdamore/tcell/v2"

type StripColor struct {
	index  int
	colors []tcell.Color
}

func NewStripColor(colors ...tcell.Color) *StripColor {
	if len(colors) < 2 {
		panic("minimum 2 arguments required")
	}
	return &StripColor{
		index:  0,
		colors: colors,
	}
}

func (s *StripColor) Next() tcell.Color {
	color := s.colors[s.index]
	if s.index == len(s.colors)-1 {
		s.index = 0
	} else {
		s.index++
	}
	return color
}

func (s *StripColor) Flush() {
	s.index = 0
}
