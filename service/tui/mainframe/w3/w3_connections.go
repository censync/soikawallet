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
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageW3Connections struct {
	*twoframes.BaseFrame
	*state.State

	// vars
	connections map[string]string
}

func NewPageW3Connections(state *state.State) *pageW3Connections {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(1, 1, 10, 10)

	return &pageW3Connections{
		State:       state,
		BaseFrame:   twoframes.NewBaseFrame(layout),
		connections: map[string]string{},
	}
}

func (p *pageW3Connections) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.EmitW3(events.EventW3InternalGetConnections, nil)
	} else {
		connections, ok := p.Params()[0].(map[string]string)
		if !ok {
			p.Emit(events.EventLogError, "Cannot parse connections request")
			p.SwitchToPage(p.Pages().GetPrevious())
		}

		for connectionId, connectionInfo := range connections {
			btnDisconnectEntry := tview.NewButton("Disconnect")

			connectionEntryFormat := fmt.Sprintf("%s - %s", connectionId, connectionInfo)

			connectionEntry := tview.NewFlex().
				SetDirection(tview.FlexColumn).
				AddItem(tview.NewTextView().SetText(connectionEntryFormat), 0, 5, false).
				AddItem(nil, 0, 1, false).
				AddItem(btnDisconnectEntry, 0, 1, false).
				AddItem(nil, 0, 1, false)

			connectionEntry.SetBorder(true)

			p.BaseLayout().AddItem(connectionEntry, 3, 1, false)
		}
		// TODO: Add reload pages
	}

}
