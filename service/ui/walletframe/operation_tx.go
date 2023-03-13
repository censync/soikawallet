package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/ui/handler"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/formtextview"
	"github.com/rivo/tview"
)

type pageOperationTx struct {
	*BaseFrame
	*state.State

	selectedAddr *responses.AddressResponse

	// ui
	layoutTokensTreeView *tview.TreeView
	layoutOperationForm  *tview.Form

	// vars
	availableTokens *responses.AddressTokensListResponse
	tokensList      []string
}

func newPageOperationTx(state *state.State) *pageOperationTx {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageOperationTx{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageOperationTx) FuncOnShow() {
	var err error
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			handler.EventLogError,
			fmt.Sprintf("Sender address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.selectedAddr = p.Params()[0].(*responses.AddressResponse)

	p.availableTokens, err = p.API().GetTokensByPath(&dto.GetAddressTokensByPathDTO{
		DerivationPath: p.selectedAddr.Path,
	})

	if err != nil {
		p.Emit(
			handler.EventLogError,
			fmt.Sprintf("Cannot get available tokens"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.tokensList = make([]string, 0)

	for _, token := range *p.availableTokens {
		p.tokensList = append(p.tokensList, token.Symbol)
	}
	p.tokensList = append(p.tokensList, "Add token")

	p.layoutTokensTreeView = tview.NewTreeView()
	p.layoutOperationForm = p.uiOperationForm()
	layoutOperation := tview.NewFlex().
		SetDirection(tview.FlexColumn)
	layoutOperation.SetBorder(true)
	layoutOperation.AddItem(p.layoutTokensTreeView, 0, 1, false)
	layoutOperation.AddItem(p.layoutOperationForm, 0, 2, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(layoutOperation, 0, 4, false).
		AddItem(nil, 0, 1, false)

	go p.actionUpdateTokens()
}

func (p *pageOperationTx) actionUpdateTokens() {
	nodeTokens := tview.NewTreeNode("tokens")
	p.layoutTokensTreeView.SetRoot(nodeTokens).
		SetTopLevel(1)
	p.layoutTokensTreeView.SetBorder(true)

	balances, err := p.API().GetTokensBalancesByPath(&dto.GetAddressTokensByPathDTO{
		DerivationPath: p.selectedAddr.Path,
	})

	if err != nil {
		p.Emit(
			handler.EventLogError,
			fmt.Sprintf("Cannot get data for %s: %s", p.selectedAddr.Path, err),
		)
	}

	for key, value := range balances {
		tokenNode := tview.NewTreeNode(fmt.Sprintf("$%s - %f", key, value))
		nodeTokens.AddChild(tokenNode)
	}

	p.Emit(handler.EventDrawForce, nil)
}

func (p *pageOperationTx) uiOperationForm() *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	inputAddrSender := formtextview.NewFormTextView(p.selectedAddr.Address)

	inputAddrReceiver := tview.NewInputField().
		SetLabel(`Receiver`)

	inputAddrAmount := tview.NewInputField().
		SetLabel(`Amount`)

	inputAddrCurrency := tview.NewDropDown().
		SetLabel("Currency").
		SetFieldWidth(10).
		SetOptions(p.tokensList, func(text string, index int) {
			if index == len(p.tokensList)-1 {
				p.SwitchToPage(pageNameTokenAdd, p.selectedAddr.CoinType, p.selectedAddr.Path)
			}
		}).
		SetCurrentOption(0)

	layoutForm.AddFormItem(inputAddrSender).
		AddFormItem(inputAddrReceiver).
		AddFormItem(inputAddrAmount).
		AddFormItem(inputAddrCurrency).
		AddButton("Send", func() {

		})
	return layoutForm
}

func (p *pageOperationTx) uiConfirmSendForm(receiver, contract string) *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	/* inputAddrSender := formtextview.NewFormTextView(p.selectedAddr.Address)

	inputAddrReceiver := formtextview.NewFormTextView(fmt.Sprintf(`Receiver: %s`, receiver))

	p.availableTokens[currency]

	inputAddrCurrency := formtextview.NewFormTextView(fmt.Sprintf(`%s Max: %f`, currency, 11111.222))
	inputAddrAmount := tview.NewInputField().
		SetLabel(`Amount`)

	layoutForm.AddFormItem(inputAddrSender).
		AddFormItem(inputAddrReceiver).
		AddFormItem(inputAddrAmount).
		AddFormItem(inputAddrCurrency).
		AddButton("Send", func() {

		}) */
	return layoutForm
}

func (p *pageOperationTx) FuncOnHide() {
	p.selectedAddr = nil
	p.tokensList = nil
	p.availableTokens = nil
	p.layout.Clear()
}
