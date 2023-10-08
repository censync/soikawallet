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
// along with  soikawallet. If not, see <http://www.gnu.org/licenses/>.

package w3

import (
	"fmt"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageW3RequestAccounts struct {
	*twidget.BaseFrame
	*state.State
}

func NewPageW3RequestAccounts(state *state.State) *pageW3RequestAccounts {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3RequestAccounts{
		State:     state,
		BaseFrame: twidget.NewBaseFrame(layout),
	}
}

func (p *pageW3RequestAccounts) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Request address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}
	//connectionReq := p.Params()[0].(*dto.ConnectDTO)

}
