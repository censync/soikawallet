package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/rivo/tview"
)

type pageTransactions struct {
	*BaseFrame
	*state.State

	//
	selectedTx      string
	selectedNetwork uint32
}

func newPageTransactions(state *state.State) *pageTransactions {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	layout.SetBorderPadding(1, 0, 0, 0)

	return &pageTransactions{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageTransactions) FuncOnShow() {

	layoutReceipt := tview.NewFlex().
		SetDirection(tview.FlexRow)
	layoutReceipt.SetBorder(true)

	searchForm := tview.NewForm().
		SetHorizontal(true)

	receiptView := tview.NewTextView()

	searchForm.
		AddInputField("Hash", p.selectedTx, 65, nil, nil).
		AddButton("search", func() {
			if p.API() != nil {
				receiptView.Clear()
				receipt, err := p.API().GetTxReceipt(&dto.GetTxReceiptDTO{
					DerivationPath: "m/44'/60'/1'/0/0'",
					Hash:           searchForm.GetFormItem(0).(*tview.InputField).GetText(),
				})
				if err == nil {
					str := ""
					for key, value := range receipt {
						str += fmt.Sprintf("%s: [darkcyan]%s\n", key, value)
					}
					receiptView.SetText(str)
				} else {
					p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get tx receipt: %s", err))
				}
			}
		})

	layoutReceipt.AddItem(searchForm, 3, 1, false).
		AddItem(receiptView, 0, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutReceipt, 0, 4, false).
		AddItem(nil, 0, 1, false)

}

func (p *pageTransactions) FuncOnHide() {
	p.selectedTx = ``
	p.layout.Clear()
}
