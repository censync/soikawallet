package formtextview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
