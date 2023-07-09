package flexmenu

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FlexMenu struct {
	largeButtons bool
	*tview.Flex
	items []*menuItem
}

func NewFlexMenu(largeButtons bool) *FlexMenu {
	menuLayout := tview.NewFlex().
		SetDirection(tview.FlexRow)

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
