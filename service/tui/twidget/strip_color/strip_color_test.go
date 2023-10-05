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
