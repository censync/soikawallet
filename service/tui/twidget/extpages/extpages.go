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

package extpages

import (
	"github.com/censync/tview"
	"github.com/gdamore/tcell/v2"
)

type ExtPage struct {
	name       string
	item       tview.Primitive
	params     []interface{}
	resize     bool
	visible    bool
	funcOnShow func()
	funcOnHide func()
	funcOnDraw func()
}

type ExtPages struct {
	*tview.Box

	pages []*ExtPage

	setFocus func(p tview.Primitive)

	changed func()

	previous string
	current  *ExtPage
}

func NewPage(name string, item tview.Primitive, resize, visible bool, onShow, onHide, onDraw func()) *ExtPage {
	return &ExtPage{
		name:       name,
		item:       item,
		resize:     resize,
		visible:    visible,
		funcOnShow: onShow,
		funcOnHide: onHide,
		funcOnDraw: onDraw,
	}
}

func (p *ExtPage) Name() string {
	return p.name
}

func (p *ExtPage) Item() tview.Primitive {
	return p.item
}

func (p *ExtPage) Resize() bool {
	return p.resize
}

func (p *ExtPage) Visible() bool {
	return p.visible
}

func (p *ExtPage) Params() []interface{} {
	return p.params
}

func (p *ExtPage) SetVisible(visible bool) {
	if p.visible != visible {
		p.visible = visible
		if p.visible {
			if p.funcOnShow != nil {
				p.funcOnShow()
			}
		} else {
			if p.funcOnHide != nil {
				p.funcOnHide()
			}
		}
	}
}

func (p *ExtPage) SetFuncOnShow(handler func()) *ExtPage {
	p.funcOnShow = handler
	return p
}

func (p *ExtPage) SetFuncOnHide(handler func()) *ExtPage {
	p.funcOnHide = handler
	return p
}

func (p *ExtPage) SetFuncOnDraw(handler func()) *ExtPage {
	p.funcOnDraw = handler
	return p
}

// NewPages returns a new ExtPages object.
func NewPages() *ExtPages {
	p := &ExtPages{
		Box: tview.NewBox(),
	}
	return p
}

func (p *ExtPages) SetChangedFunc(handler func()) *ExtPages {
	p.changed = handler
	return p
}

func (p *ExtPages) GetPageCount() int {
	return len(p.pages)
}

func (p *ExtPages) AddPage(page *ExtPage) *ExtPages {
	if page == nil {
		panic("cannot add not initialised pages")
	}
	hasFocus := p.HasFocus()
	for index, pg := range p.pages {
		if pg.Name() == page.Name() {
			p.pages = append(p.pages[:index], p.pages[index+1:]...)
			break
		}
	}
	p.pages = append(p.pages, page)
	if p.changed != nil {
		p.changed()
	}
	if hasFocus {
		p.Focus(p.setFocus)
	}
	return p
}

func (p *ExtPages) RemovePage(name string) *ExtPages {
	var isVisible bool
	hasFocus := p.HasFocus()
	if p.previous == name {
		p.previous = ``
	}
	for index, page := range p.pages {
		if page.Name() == name {
			isVisible = page.Visible()
			p.pages = append(p.pages[:index], p.pages[index+1:]...)
			if page.Visible() && p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if isVisible {
		for index, page := range p.pages {
			if index < len(p.pages)-1 {
				if page.Visible() {
					break // There is a remaining visible pages.
				}
			} else {
				p.current = page
				page.SetVisible(true)
			}
		}
	}
	if hasFocus {
		p.Focus(p.setFocus)
	}
	return p
}

func (p *ExtPages) GetPrevious() string {
	return p.previous
}

func (p *ExtPages) IsPreviousExists() bool {
	return p.previous != ""
}

func (p *ExtPages) HasPage(name string) bool {
	for _, page := range p.pages {
		if page.Name() == name {
			return true
		}
	}
	return false
}

func (p *ExtPages) HidePage(name string) *ExtPages {
	for _, page := range p.pages {
		if page.Name() == name {
			page.SetVisible(false)
			if p.changed != nil {
				p.changed()
			}
			break
		}
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

func (p *ExtPages) SwitchToPage(name string, params ...interface{}) *ExtPages {
	if p.current != nil && p.current.name != name {
		p.previous = p.current.name
	}
	for id, page := range p.pages {
		if page.Name() == name {
			if len(params) > 0 {
				// TODO: Check for leaks
				if p.pages[id].params != nil {
					p.pages[id].params = nil
				}
				p.pages[id].params = make([]interface{}, len(params))
				copy(p.pages[id].params, params)
			}
			p.current = page
			page.SetVisible(true)
		} else {
			if page.Visible() {
				page.SetVisible(false)
			}
		}
	}
	if p.changed != nil {
		p.changed()
	}
	if p.HasFocus() {
		p.Focus(p.setFocus)
	}
	return p
}

func (p *ExtPages) GetFrontPage() (name string, item tview.Primitive) {
	for index := len(p.pages) - 1; index >= 0; index-- {
		if p.pages[index].Visible() {
			return p.pages[index].Name(), p.pages[index].Item()
		}
	}
	return
}

func (p *ExtPages) HasFocus() bool {
	for _, page := range p.pages {
		if page.Item().HasFocus() {
			return true
		}
	}
	return p.Box.HasFocus()
}

func (p *ExtPages) Focus(delegate func(p tview.Primitive)) {
	if delegate == nil {
		return // We cannot delegate so we cannot focus.
	}
	p.setFocus = delegate
	var topItem tview.Primitive
	for _, page := range p.pages {
		if page.Visible() {
			topItem = page.Item()
		}
	}
	if topItem != nil {
		delegate(topItem)
	} else {
		p.Box.Focus(delegate)
	}
}

func (p *ExtPages) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)
	for _, page := range p.pages {
		if !page.Visible() {
			continue
		}
		if page.Resize() {
			x, y, width, height := p.GetInnerRect()
			page.Item().SetRect(x, y, width, height)
		}
		page.Item().Draw(screen)
		if page.funcOnDraw != nil {
			page.funcOnDraw()
		}
	}
}

func (p *ExtPages) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return p.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !p.InRect(event.Position()) {
			return false, nil
		}

		for index := len(p.pages) - 1; index >= 0; index-- {
			page := p.pages[index]
			if page.Visible() {
				consumed, capture = page.Item().MouseHandler()(action, event, setFocus)
				if consumed {
					return
				}
			}
		}

		return
	})
}

func (p *ExtPages) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		for _, page := range p.pages {
			if page.Item().HasFocus() {
				if handler := page.Item().InputHandler(); handler != nil {
					handler(event, setFocus)
					return
				}
			}
		}
	})
}

func (p *ExtPages) Current() *ExtPage {
	return p.current
}
