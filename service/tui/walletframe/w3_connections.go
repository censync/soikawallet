package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/rivo/tview"
)

type pageW3Connections struct {
	*BaseFrame
	*state.State

	// vars
	connections map[string]string
}

func newPageW3Connections(state *state.State) *pageW3Connections {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layout.SetBorderPadding(10, 10, 10, 10)

	return &pageW3Connections{
		State:       state,
		BaseFrame:   &BaseFrame{layout: layout},
		connections: map[string]string{},
	}
}

func (p *pageW3Connections) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 1 {
		p.EmitW3(event_bus.EventW3InternalGetConnections, nil)
	} else {
		connections, ok := p.Params()[0].(map[string]string)
		if !ok {
			p.Emit(event_bus.EventLogError, "Cannot parse connections request")
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

			p.layout.AddItem(connectionEntry, 5, 1, false)
		}
		// TODO: Add reload page
	}

}

func (p *pageW3Connections) FuncOnHide() {
	p.layout.Clear()
}
