package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/widgets/formtextview"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
)

type pageTokenAdd struct {
	*BaseFrame
	*state.State

	// ui
	layoutTokenAdd *tview.Flex

	// vars
	paramSelectedNetwork        types.NetworkType
	selectedTokenStandard       types.TokenStandard
	paramSelectedDerivationPath string
}

func newPageTokenAdd(state *state.State) *pageTokenAdd {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageTokenAdd{
		State:     state,
		BaseFrame: &BaseFrame{layout: layout},
	}
}

func (p *pageTokenAdd) Layout() *tview.Flex {
	p.layoutTokenAdd = tview.NewFlex().
		SetDirection(tview.FlexRow)
	p.layoutTokenAdd.SetBorder(true)
	return p.layout
}

func (p *pageTokenAdd) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 2 {
		p.Emit(
			event_bus.EventLogError,
			fmt.Sprintf("Incorrect params"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.paramSelectedNetwork = p.Params()[0].(types.NetworkType)
	p.paramSelectedDerivationPath = p.Params()[1].(string)

	p.layoutTokenAdd.AddItem(p.uiTokenAddForm(), 0, 1, false)

	p.layout.AddItem(nil, 0, 1, false).
		AddItem(p.layoutTokenAdd, 0, 4, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageTokenAdd) uiTokenAddForm() *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	inputContractAddr := tview.NewInputField().
		SetLabel(`Contract address`).
		SetText("0x91B268bd44c6a16b2E518060b44eFF33cB17f84d") // debug

	inputSelectTokenStandard := tview.NewDropDown().
		SetLabel("Type").
		SetFieldWidth(10).
		SetOptions(types.GetTokenStandardNames(p.paramSelectedNetwork), func(text string, index int) {
			p.selectedTokenStandard = types.GetTokenStandByName(text)
		}).
		SetCurrentOption(0)

	layoutForm.
		AddFormItem(inputContractAddr).
		AddFormItem(inputSelectTokenStandard).
		AddButton("Check contract", func() {
			tokenInfo, err := p.API().GetToken(&dto.GetTokenDTO{
				NetworkType: uint32(p.paramSelectedNetwork),
				Contract:    inputContractAddr.GetText(),
			})
			if err != nil {
				p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get token data: %s", err))
			} else {
				p.layoutTokenAdd.Clear()
				p.layoutTokenAdd.AddItem(p.uiTokenConfirmForm(tokenInfo), 0, 1, false)
			}
		})
	return layoutForm
}

func (p *pageTokenAdd) uiTokenConfirmForm(tokenConfig *resp.TokenConfig) *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	inputContractAddr := formtextview.NewFormTextView(
		fmt.Sprintf("[yellow]Contract: [white]%s\n[yellow]Name: [white]%s\n[yellow]Symbol: [white]%s\n[yellow]Decimals: [white]%d",
			tokenConfig.Contract,
			tokenConfig.Name,
			tokenConfig.Symbol,
			tokenConfig.Decimals,
		),
	)

	layoutForm.AddFormItem(inputContractAddr).
		AddButton("Add token", func() {
			err := p.API().UpsertToken(&dto.AddTokenDTO{
				Standard:       uint8(p.selectedTokenStandard),
				NetworkType:    uint32(p.paramSelectedNetwork),
				Contract:       tokenConfig.Contract,
				DerivationPath: p.paramSelectedDerivationPath,
			})

			if err != nil {
				p.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot add token: %s", err))
			} else {
				p.Emit(event_bus.EventLogSuccess, fmt.Sprintf(
					"Added token \"%s\" - \"%s\"",
					tokenConfig.Name,
					tokenConfig.Symbol,
				),
				)
				p.layoutTokenAdd.Clear()
				p.SwitchToPage(p.Pages().GetPrevious())
			}
		})
	return layoutForm
}

func (p *pageTokenAdd) FuncOnHide() {
	p.layout.Clear()
}
