package walletframe

import (
	"encoding/json"
	"fmt"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/rivo/tview"
	"os"
)

func (p *pageSettings) tabApp() *tview.Flex {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)

	formConfigLocation := tview.NewForm().
		SetHorizontal(true).
		AddInputField("path", ".soikawallet/config", 20, nil, nil).
		AddButton("Save", func() {

		})
	formSaveConfig := tview.NewForm().
		SetHorizontal(true).
		AddButton("Export", func() {
			// file, err := os.Open("config")
			data, err := json.Marshal(p.API())
			if err != nil {
				p.Emit(handler.EventLogError, fmt.Sprintf("Cannot marshal config: %s", err))
				return
			}
			err = os.WriteFile("config.json", data, 644)
			if err != nil {
				p.Emit(handler.EventLogError, fmt.Sprintf("Cannot save file: %s", err))
				return
			}
		})

	layout.AddItem(formConfigLocation, 3, 1, false).
		AddItem(formSaveConfig, 3, 1, false)
	return layout
}
