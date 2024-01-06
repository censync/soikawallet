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

package token

import (
	"fmt"
	mhda "github.com/censync/go-mhda"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/tview"
	"github.com/censync/twidget/twoframes"
)

type pageTokenAdd struct {
	*twoframes.BaseFrame
	*state.State

	// ui
	layoutTokenAdd *tview.Flex

	// vars
	paramSelectedChainKey       *mhda.Chain
	selectedTokenStandard       string
	paramSelectedDerivationPath string
}

func NewPageTokenAdd(state *state.State) *pageTokenAdd {
	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn)

	return &pageTokenAdd{
		State:     state,
		BaseFrame: twoframes.NewBaseFrame(layout),
	}
}

func (p *pageTokenAdd) Layout() *tview.Flex {
	p.layoutTokenAdd = tview.NewFlex().
		SetDirection(tview.FlexRow)
	p.layoutTokenAdd.SetBorder(true)

	return p.BaseFrame.Layout()
}

func (p *pageTokenAdd) FuncOnShow() {
	if p.Params() == nil || len(p.Params()) != 2 {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Incorrect params"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	// TODO: Add exception handling
	chainKey, err := mhda.ChainFromKey(p.Params()[0].(mhda.ChainKey))

	if err != nil {
		p.Emit(
			events.EventLogError,
			fmt.Sprintf("Incorrect chainKey"),
		)
		p.SwitchToPage(p.Pages().GetPrevious())
	}

	p.paramSelectedChainKey = chainKey

	p.paramSelectedDerivationPath = p.Params()[1].(string)

	p.layoutTokenAdd.AddItem(p.uiTokenAddForm(), 0, 1, false)

	p.BaseLayout().AddItem(nil, 0, 1, false).
		AddItem(p.layoutTokenAdd, 0, 4, false).
		AddItem(nil, 0, 1, false)
}

func (p *pageTokenAdd) uiTokenAddForm() *tview.Form {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	inputContractAddr := tview.NewInputField().
		SetLabel(`Contract address`).
		SetText("0x8A953CfE442c5E8855cc6c61b1293FA648BAE472") // debug

	inputSelectTokenStandard := tview.NewDropDown().
		SetLabel("Type").
		SetFieldWidth(10).
		SetOptions(p.API().GetTokenStandardNamesByChain(&dto.GetTokenStandardNamesByNetworkDTO{
			NetworkType: p.paramSelectedChainKey.NetworkType(),
		}), func(text string, index int) {
			p.selectedTokenStandard = text
		}).
		SetCurrentOption(0)

	layoutForm.
		AddFormItem(inputContractAddr).
		AddFormItem(inputSelectTokenStandard).
		AddButton("Check contract", func() {
			tokenInfo, err := p.API().GetToken(&dto.GetTokenDTO{
				ChainKey: p.paramSelectedChainKey.Key(),
				Contract: inputContractAddr.GetText(),
			})
			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot get token data: %s", err))
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

	inputContractAddr := tview.NewTextView().SetText(
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
				Standard: p.selectedTokenStandard,
				ChainKey: p.paramSelectedChainKey.Key(),
				Contract: tokenConfig.Contract,
				MhdaPath: p.paramSelectedDerivationPath,
			})

			if err != nil {
				p.Emit(events.EventLogError, fmt.Sprintf("Cannot add token: %s", err))
			} else {
				p.Emit(events.EventLogSuccess, fmt.Sprintf(
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
