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
