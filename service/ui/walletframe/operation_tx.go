package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/formtextview"
	"github.com/rivo/tview"
	"strconv"
)

type pageOperationTx struct {
	*BaseFrame
	*state.State

	paramSelectedAddr *responses.AddressResponse

	// ui
	layoutTokensTreeView *tview.TreeView
	layoutOperationForm  *tview.Form

	// vars
	availableTokens *responses.AddressTokensListResponse
	tokensList      []string
	tokensMap       map[int]string

	selectedToken *responses.AddressTokenEntry
}

func newPageOperationTx(state *state.State) *pageOperationTx {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageOperationTx{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
		tokensMap: map[int]string{},
	}
}

func (p *pageOperationTx) FuncOnShow() {
	var err error
	if p.Params() == nil || len(p.Params()) != 1 {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Sender address is not set"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.paramSelectedAddr = p.Params()[0].(*responses.AddressResponse)

	p.availableTokens, err = p.API().GetTokensByPath(&dto.GetAddressTokensByPathDTO{
		DerivationPath: p.paramSelectedAddr.Path,
	})

	if err != nil {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Cannot get available tokens"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.tokensList = make([]string, 0)

	index := 0
	for contract, token := range *p.availableTokens {
		p.tokensList = append(p.tokensList, token.Symbol)
		p.tokensMap[index] = contract
		index++
	}
	p.tokensList = append(p.tokensList, " [ add token] ")

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
		DerivationPath: p.paramSelectedAddr.Path,
	})

	if err != nil {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Cannot get data for %s: %s", p.paramSelectedAddr.Path, err),
		)
	}

	for key, value := range balances {
		tokenNode := tview.NewTreeNode(fmt.Sprintf("$%s - %f", key, value))
		nodeTokens.AddChild(tokenNode)
	}

	p.Emit(event_bus.EventDrawForce, nil)
}

func (p *pageOperationTx) uiOperationForm() *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	inputAddrSender := formtextview.NewFormTextView(p.paramSelectedAddr.Address)

	inputAddrReceiver := tview.NewInputField().
		SetLabel(`Receiver`)

	inputAddrAmount := tview.NewInputField().
		SetLabel(`Amount`)

	inputAddrCurrency := tview.NewDropDown().
		SetLabel("Currency").
		SetFieldWidth(10).
		SetOptions(p.tokensList, func(text string, index int) {
			if index == len(p.tokensList)-1 {
				p.SwitchToPage(pageNameTokenAdd, p.paramSelectedAddr.CoinType, p.paramSelectedAddr.Path)
			} else {
				if contract, ok := p.tokensMap[index]; ok {
					p.selectedToken = (*p.availableTokens)[contract]
				} else {
					p.Emit(event_bus.EventLogError, "Undefined token")
				}
			}
		}).
		SetCurrentOption(0)

	layoutForm.AddFormItem(inputAddrSender).
		AddFormItem(inputAddrReceiver).
		AddFormItem(inputAddrAmount).
		AddFormItem(inputAddrCurrency).
		AddButton("Send", func() {
			value, err := strconv.ParseFloat(inputAddrAmount.GetText(), 64)

			if err == nil {
				txId, err := p.API().SendTokens(&dto.SendTokensDTO{
					DerivationPath: p.paramSelectedAddr.Path,
					To:             inputAddrReceiver.GetText(),
					Value:          value,
					Standard:       p.selectedToken.Standard,
					Contract:       p.selectedToken.Contract,
				})
				if err == nil {
					p.Emit(event_bus.EventLogSuccess, fmt.Sprintf("Transaction sent: %s", txId))
				} else {
					p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot send transaction: %s", err))
				}
			} else {
				p.Emit(event_bus.EventLogError, "Incorrect value")
			}
		})
	return layoutForm
}

func (p *pageOperationTx) uiConfirmSendForm(receiver, contract string) *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	/* inputAddrSender := formtextview.NewFormTextView(p.paramSelectedAddr.Address)

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
	p.paramSelectedAddr = nil
	p.tokensList = nil
	p.availableTokens = nil
	p.layout.Clear()
}
