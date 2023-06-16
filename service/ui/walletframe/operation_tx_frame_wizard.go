package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/ui/state"
	"github.com/censync/soikawallet/service/ui/widgets/formtextview"
	"github.com/censync/soikawallet/types"
	"github.com/rivo/tview"
	"strconv"
)

type frameOperationWizard struct {
	layout *tview.Flex
	*state.State

	// vars
	selectedAddress   *resp.AddressResponse
	selectedToken     *resp.AddressTokenEntry
	selectedRecipient string
	selectedAmount    float64
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
		SetText("0x8680C23520C3731024a0632277FcB7303cD1A00f")

	inputValue := tview.NewInputField().
		SetAcceptanceFunc(tview.InputFieldFloat).
		SetLabel(`Amount`)

	tokensList := make([]string, 0)
	tokensMap := map[int]string{}

	availableTokens, err := f.API().GetTokensByPath(&dto.GetAddressTokensByPathDTO{
		DerivationPath: f.selectedAddress.Path,
	})

	if err != nil {
		f.Emit(event_bus.EventLogError, "Cannot get available tokens")
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
				f.SwitchToPage(pageNameTokenAdd, f.selectedAddress.CoinType, f.selectedAddress.Path)
			} else {
				if contract, ok := tokensMap[index]; ok {
					f.selectedToken = (*availableTokens)[contract]
				} else {
					f.Emit(event_bus.EventLogError, "Undefined token")
				}
			}
		}).
		SetCurrentOption(0)

	layoutForm.AddFormItem(inputAddrSender).
		AddFormItem(inputAddrReceiver).
		AddFormItem(inputValue).
		AddFormItem(inputAddrCurrency).
		AddButton("Send", func() {
			f.selectedAmount, err = strconv.ParseFloat(inputValue.GetText(), 64)
			if err != nil {
				f.Emit(event_bus.EventLogError, "Incorrect value")
				return
			}
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
	if types.TokenStandard(f.selectedToken.Standard) != types.TokenBase {
		allowance, err := f.API().GetAllowance(&dto.GetTokenAllowanceDTO{
			DerivationPath: f.selectedAddress.Path,
			To:             f.selectedRecipient,
			Value:          f.selectedAmount,
			Standard:       f.selectedToken.Standard,
			Contract:       f.selectedToken.Contract,
		})

		if err != nil {
			f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get allowance: %e", err))
			return false
		}

		if allowance == 0 {
			f.Emit(event_bus.EventLogWarning, "Not approved, zero allowance")
			return false
		}
	}
	return true
}

func (f *frameOperationWizard) actionConfigureGas() {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	layoutForm.AddTextView("title", "actionConfigureGas", 30, 1, true, false)

	gasConfig, err := f.API().GetGasPriceBaseTx(&dto.GetGasPriceBaseTxDTO{
		DerivationPath: f.selectedAddress.Path,
	})

	if err == nil || gasConfig == nil {
		f.Emit(event_bus.EventLogInfo, gasConfig)

		/*
			gasCalc := types.NewGasCalculator(types.GasAlgEVML1v1, &types.GasCalcOpts{
				GasSuffix:    "gwei",
				TokenSuffix:  "Eth",
				FiatSuffix:   "USD",
				FiatCurrency: 1700,
			})
		*/

		calculatedFee := gasConfig["units"] * (gasConfig["base_fee"] + gasConfig["priority_fee"]) / 1e9

		labelCalcFee := tview.NewTextView().
			SetText(fmt.Sprintf("Total fee: %f", calculatedFee))

		inputBaseFee := tview.NewInputField().
			SetFieldWidth(10).
			SetLabel("Base fee").
			SetText(fmt.Sprintf("%f", gasConfig["base_fee"]))

		inputPriorityFee := tview.NewInputField().
			SetFieldWidth(10).
			SetLabel("Priority fee"). // 0.1 ../ 0.25 / 3 / 4 / .. 5 (ETH)
			SetText(fmt.Sprintf("%f", gasConfig["priority_fee"]))

		layoutForm.
			AddFormItem(inputBaseFee).
			AddFormItem(inputPriorityFee).
			AddFormItem(labelCalcFee).
			AddButton("Send", func() {

			})
	} else {
		// add label
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot send transaction: %s", err))
	}

	f.layout.Clear()
	f.layout.AddItem(layoutForm, 0, 1, false)
	//f.Emit(event_bus.EventDrawForce, nil)
}

func (f *frameOperationWizard) actionConfigureAllowance() {
	layoutForm := tview.NewForm().
		SetHorizontal(false)
	layoutForm.SetBorder(true)

	layoutForm.AddTextView("title", "actionConfigureAllowance", 30, 1, true, false)
	f.layout.Clear()
	f.layout.AddItem(layoutForm, 0, 1, false)
	//f.Emit(event_bus.EventDrawForce, nil)
}

func (f *frameOperationWizard) actionSendTransaction() {
	txId, err := f.API().SendTokens(&dto.SendTokensDTO{
		DerivationPath: f.selectedAddress.Path,
		To:             f.selectedRecipient,
		Value:          f.selectedAmount,
		Standard:       f.selectedToken.Standard,
		Contract:       f.selectedToken.Contract,
	})
	if err == nil {
		f.Emit(event_bus.EventLogSuccess, fmt.Sprintf("Transaction sent: %s", txId))
	} else {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot send transaction: %s", err))
	}
}
