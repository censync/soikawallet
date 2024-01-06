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
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package op_tx

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/tui/events"
	"github.com/censync/soikawallet/service/tui/pages"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/types/gas"
	"github.com/censync/tview"
	"github.com/censync/twidget/formtextview"
)

type frameOperationWizard struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedAddress   *resp.AddressResponse
	selectedToken     *resp.AddressTokenEntry
	selectedRecipient string
	selectedAmount    string
}

func newFrameOperationWizard(state *state.State, selectedAddress *resp.AddressResponse) *frameOperationWizard {
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow)
	return &frameOperationWizard{
		State:           state,
		layout:          layout,
		selectedAddress: selectedAddress,
	}
}

func (f *frameOperationWizard) Layout() *tview.Flex {
	var err error

	layoutForm := tview.NewForm().
		SetHorizontal(false)

	inputAddrSender := formtextview.NewFormTextView(f.selectedAddress.Address)

	inputAddrReceiver := tview.NewInputField().
		SetLabel(`Receiver`).
		SetText("")

	inputValue := tview.NewInputField().
		SetAcceptanceFunc(tview.InputFieldFloat).
		SetLabel(`Amount`).
		SetText("0.001")

	tokensList := make([]string, 0)
	tokensMap := map[int]string{}

	availableTokens, err := f.API().GetTokensByPath(&dto.GetAddressTokensByPathDTO{
		MhdaPath: f.selectedAddress.Path,
	})

	if err != nil {
		f.Emit(events.EventLogError, "Cannot get available tokens")
		f.SwitchToPage(f.Pages().GetPrevious())
	}

	index := 0
	for contract, token := range *availableTokens {
		tokensList = append(tokensList, token.Symbol)
		tokensMap[index] = contract
		index++
	}
	tokensList = append(tokensList, " [ add token] ")

	inputAddrCurrency := tview.NewDropDown().
		SetLabel("Currency").
		SetFieldWidth(10).
		SetOptions(tokensList, func(text string, index int) {
			if index == len(tokensList)-1 {
				f.SwitchToPage(pages.TokenAdd, f.selectedAddress.ChainKey, f.selectedAddress.Path)
			} else {
				if contract, ok := tokensMap[index]; ok {
					f.selectedToken = (*availableTokens)[contract]
				} else {
					f.Emit(events.EventLogError, "Undefined token")
				}
			}
		}).
		SetCurrentOption(0)

	layoutForm.AddFormItem(inputAddrSender).
		AddFormItem(inputAddrReceiver).
		AddFormItem(inputValue).
		AddFormItem(inputAddrCurrency).
		AddButton("Send", func() {
			// TODO: Add float validation
			f.selectedAmount = inputValue.GetText()
			f.selectedRecipient = inputAddrReceiver.GetText()
			f.actionCheckAndStart()
		})
	f.layout.AddItem(layoutForm, 0, 1, false)
	return f.layout
}

func (f *frameOperationWizard) actionCheckAndStart() {
	if f.actionCheckAllowancePermission() {
		f.actionConfigureGas()
	} else {
		f.actionConfigureAllowance()
	}
}

func (f *frameOperationWizard) actionCheckAllowancePermission() bool {
	// TODO: Remove magic 1 == TokenBase
	if f.selectedToken.Standard == `Base` {
		return true
	}
	allowance, err := f.API().GetAllowance(&dto.GetTokenAllowanceDTO{
		MhdaPath: f.selectedAddress.Path,
		To:       f.selectedRecipient,
		Value:    f.selectedAmount,
		Standard: f.selectedToken.Standard,
		Contract: f.selectedToken.Contract,
	})

	if err != nil {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot get allowance: %e", err))
		return false
	}

	f.Emit(events.EventLogInfo, fmt.Sprintf("Allowance: %d ", allowance))

	// TODO: Add human check
	/*if allowance == 0 {
		f.Emit(events.EventLogWarning, "Not approved, zero allowance")
		return false
	} else {
		f.Emit(events.EventLogInfo, fmt.Sprintf("Allowance: %d ", allowance))
	}*/
	return true
}

func (f *frameOperationWizard) actionConfigureGas() {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	layoutForm.AddTextView("title", "actionConfigureGas", 30, 1, true, false)

	calcConfig, err := f.API().GetGasCalculatorConfig(&dto.GetGasCalculatorConfigDTO{
		Operation: "transfer",
		MhdaPath:  f.selectedAddress.Path,
		To:        f.selectedRecipient,
		Value:     f.selectedAmount,
		Standard:  f.selectedToken.Standard,
		Contract:  f.selectedToken.Contract,
	})

	if err != nil {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot get gas calculator instance: %s", err))
		return
	}

	gasCalc, err := gas.Unmarshal(calcConfig.Calculator)

	if err != nil {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot unmarshal gas calculator: %s", err))
		return
	}

	f.Emit(events.EventLogInfo, gasCalc.Debug())

	templateGas := fmt.Sprintf("GasEstimate: %d Base: %s [%d]\nSuggestSlow: %s [%d]\nSuggestRegular: %s [%d]\nSuggestPriority: %s [%d]\nEst gas price: %s Limit gas price: %s",
		gasCalc.EstimateGas(),
		gasCalc.FormatHumanGas(gasCalc.BaseGas()),
		gasCalc.BaseGas(),
		gasCalc.FormatHumanGas(gasCalc.SuggestSlow()),
		gasCalc.SuggestSlow(),
		gasCalc.FormatHumanGas(gasCalc.SuggestRegular()),
		gasCalc.SuggestRegular(),
		gasCalc.FormatHumanGas(gasCalc.SuggestPriority()),
		gasCalc.SuggestPriority(),
		gasCalc.FormatHumanFiatPrice(gasCalc.EstimateGas()*(gasCalc.BaseGas()+gasCalc.SuggestRegular())),
		gasCalc.FormatHumanFiatPrice(gasCalc.EstimateGas()*gasCalc.LimitMaxGasFee(gasCalc.SuggestRegular())),
	)

	labelInfo := tview.NewTextView().SetText(templateGas)

	layoutForm.
		AddFormItem(labelInfo).
		AddButton("Update gas", func() {
			f.actionConfigureGas()
		}).
		AddButton("Send", func() {
			f.actionSendTransaction(gasCalc.EstimateGas(), gasCalc.SuggestRegular(), gasCalc.LimitMaxGasFee(gasCalc.SuggestRegular()), false)
		}).
		AddButton("Send AirGap", func() {
			f.actionSendTransaction(gasCalc.EstimateGas(), gasCalc.SuggestRegular(), gasCalc.LimitMaxGasFee(gasCalc.SuggestRegular()), true)
		})

	f.layout.Clear()
	f.layout.AddItem(layoutForm, 0, 1, false)
	//f.Emit(events.EventDrawForce, nil)
}

func (f *frameOperationWizard) actionConfigureAllowance() {

	calcConfig, err := f.API().GetGasCalculatorConfig(&dto.GetGasCalculatorConfigDTO{
		Operation: "approve",
		MhdaPath:  f.selectedAddress.Path,
		To:        f.selectedRecipient,
		Value:     f.selectedAmount,
		Standard:  f.selectedToken.Standard,
		Contract:  f.selectedToken.Contract,
	})

	if err != nil {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot get gas calculator instance: %s", err))
		return
	}

	gasCalc, err := gas.Unmarshal(calcConfig.Calculator)

	if err != nil {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot unmarshal gas calculator: %s", err))
		return
	}

	templateGas := fmt.Sprintf("GasEstimate: %d Base: %s [%d]\nSuggestSlow: %s [%d]\nSuggestRegular: %s [%d]\nSuggestPriority: %s [%d]\nEst gas price: %s Limit gas price: %s",
		gasCalc.EstimateGas(),
		gasCalc.FormatHumanGas(gasCalc.BaseGas()),
		gasCalc.BaseGas(),
		gasCalc.FormatHumanGas(gasCalc.SuggestSlow()),
		gasCalc.SuggestSlow(),
		gasCalc.FormatHumanGas(gasCalc.SuggestRegular()),
		gasCalc.SuggestRegular(),
		gasCalc.FormatHumanGas(gasCalc.SuggestPriority()),
		gasCalc.SuggestPriority(),
		gasCalc.FormatHumanFiatPrice(gasCalc.EstimateGas()*(gasCalc.BaseGas()+gasCalc.SuggestRegular())),
		gasCalc.FormatHumanFiatPrice(gasCalc.EstimateGas()*gasCalc.LimitMaxGasFee(gasCalc.SuggestRegular())),
	)

	labelInfo := tview.NewTextView().SetText(templateGas)

	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)
	layoutForm.AddTextView("title", "actionConfigureAllowance", 30, 1, true, false)

	layoutForm.AddFormItem(labelInfo)
	layoutForm.AddButton("approve", func() {
		f.actionSendApprove(gasCalc.EstimateGas(), gasCalc.SuggestRegular(), gasCalc.LimitMaxGasFee(gasCalc.SuggestRegular()))
	})

	f.layout.Clear()
	f.layout.AddItem(layoutForm, 0, 1, false)
}
func (f *frameOperationWizard) actionSendApprove(gas, gasTipCap, gasFeePrice uint64) {
	txId, err := f.API().ApproveTokens(&dto.SendTokensDTO{
		MhdaPath:  f.selectedAddress.Path,
		To:        f.selectedRecipient,
		Value:     f.selectedAmount,
		Gas:       gas,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeePrice,
		Standard:  f.selectedToken.Standard,
		Contract:  f.selectedToken.Contract,
	})
	if err == nil {
		f.Emit(events.EventLogSuccess, fmt.Sprintf("Transaction approve sent: %s", txId))
		f.actionConfigureGas()
	} else {
		f.Emit(events.EventLogError, fmt.Sprintf("Cannot send approve transaction: %s", err))
	}
}

func (f *frameOperationWizard) actionSendTransaction(gas, gasTipCap, gasFeePrice uint64, isAirGap bool) {
	request := &dto.SendTokensDTO{
		MhdaPath:  f.selectedAddress.Path,
		To:        f.selectedRecipient,
		Value:     f.selectedAmount,
		Gas:       gas,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeePrice,
		Standard:  f.selectedToken.Standard,
		Contract:  f.selectedToken.Contract,
	}

	if !isAirGap {

		txId, err := f.API().SendTokens(request)

		if err == nil {
			f.Emit(events.EventLogSuccess, fmt.Sprintf("Transaction sent: %s", txId))
		} else {
			f.Emit(events.EventLogError, fmt.Sprintf("Cannot send transaction: %s", err))
		}
	} else {

		airGapData, err := f.API().SendTokensPrepare(request)

		if err == nil {
			f.Emit(events.EventLogSuccess, "Transaction prepared")
			f.SwitchToPage(pages.AirGapShow, airGapData)
		} else {
			f.Emit(events.EventLogError, fmt.Sprintf("Cannot prepare transaction: %s", err))
		}
	}
}
