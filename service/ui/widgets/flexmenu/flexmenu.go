package flexmenu

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FlexMenu struct {
	*tview.Flex
	items []*menuItem
}

func NewFlexMenu() *FlexMenu {
	return &FlexMenu{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
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
	f.items = append(f.items, &menuItem{
		Label:  label,
		Key:    key,
		Action: action,
	})

	return f
}

func (f *FlexMenu) Layout() *tview.Flex {
	f.SetBorderPadding(1, 1, 1, 1)

	for index := range f.items {
		if f.items[index] == nil {
			panic(fmt.Sprintf("target not set for button \"%s\"", f.items[index].Label))
		}

		btn := tview.NewButton(f.items[index].LabelDecorated()).
			SetSelectedFunc(f.items[index].Action).
			SetLabelAlign(tview.AlignLeft).
			SetActivatedStyleAttrs(tcell.AttrBold)

		btn.SetBorderPadding(0, 0, 2, 0)

		f.Flex.AddItem(btn, 1, 1, false)
		f.Flex.AddItem(nil, 1, 1, false)
	}
	return f.Flex
}
