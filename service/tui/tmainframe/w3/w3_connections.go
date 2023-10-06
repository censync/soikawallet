package w3

import (
	"fmt"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/twidget"
	"github.com/censync/tview"
)

type pageW3Connections struct {
	*twidget.BaseFrame
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
		BaseFrame:   twidget.NewBaseFrame(layout),
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
		// TODO: Add reload page
	}

}
