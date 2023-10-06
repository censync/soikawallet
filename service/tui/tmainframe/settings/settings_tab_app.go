package settings

import (
	"fmt"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/tview"
	"os"
)

func (p *pageSettings) tabApp() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	formConfigLocation := tview.NewForm().
		SetHorizontal(true).
		AddInputField("path", ".soikawallet/config.json", 20, nil, nil).
		AddButton("Export", func() {
			// DEBUG
			homeDir, err := os.UserHomeDir()
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot get user home dir: %s", err))
				return
			}

			meta, err := p.API().ExportMetaDebug()
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot get meta: %s", err))
				return
			}

			if _, err = os.Stat(homeDir + "/.soikawallet"); err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(homeDir+"/.soikawallet", os.ModePerm)
					if err != nil {
						p.Emit(events.EventLogError, fmt.Sprintf("Cannot create settings dir: %s", err))
						return
					}
				}
				p.Emit(events.EventLogError, fmt.Sprintf("Err: %s", err))
			}
			err = os.WriteFile(homeDir+"/.soikawallet/config.json", meta, 600)
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot save file: %s", err))
				return
			}
			p.Emit(events.EventLogSuccess, "Config saved")
		})

	layout.AddItem(formConfigLocation, 3, 1, false)
	return layout
}
