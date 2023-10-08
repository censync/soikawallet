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
