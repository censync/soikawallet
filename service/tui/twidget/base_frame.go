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
