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
