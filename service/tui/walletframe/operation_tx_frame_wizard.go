package walletframe

import (
	"fmt"
	"github.com/censync/soikawallet/api/dto"
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/service/internal/event_bus"
	"github.com/censync/soikawallet/service/tui/state"
	"github.com/censync/soikawallet/service/tui/widgets/formtextview"
	"github.com/censync/soikawallet/types"
	"github.com/censync/soikawallet/types/gas"
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
		SetLabel(`Amount`).
		SetText("0.001")

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
				f.SwitchToPage(pageNameTokenAdd, f.selectedAddress.NetworkType, f.selectedAddress.Path)
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

	/*gasConfig, err := f.API().GetGasPriceBaseTx(&dto.GetGasPriceBaseTxDTO{
		DerivationPath: f.selectedAddress.Path,
	})

	if err == nil || gasConfig == nil {
		f.Emit(event_bus.EventLogInfo, gasConfig)

		fiatCurrency, suffix, _ := f.API().GetFiatCurrency(&dto.GetFiatCurrencyDTO{
			NetworkType: uint32(f.selectedAddress.NetworkType),
		})

		gasCalc := gas.NewCalcEVML1V1(&gas.CalcEVML1V1{
			CalcOpts: &gas.CalcOpts{
				GasSymbol:    "gwei",
				GasUnits:     10e9,
				TokenSuffix:  suffix,
				FiatCurrency: fiatCurrency,
			},
			Units:       gasConfig["units"],
			BaseFee:     gasConfig["base_fee"],
			PriorityFee: gasConfig["priority_fee"],
			GasLimit:      30000, // 30e6?
		})

		calculatedFee := gasConfig["units"] * (gasConfig["base_fee"] + gasConfig["priority_fee"]) / 1e9
	*/

	calcConfig, err := f.API().GetGasCalculatorConfig(&dto.GetAddressCalculatorConfigDTO{
		DerivationPath: f.selectedAddress.Path,
	})

	if err != nil {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot get gas calculator instance: %s", err))
	}

	gasCalc, err := gas.Unmarshal(calcConfig.Calculator)

	if err != nil {
		f.Emit(event_bus.EventLogError, fmt.Sprintf("Cannot unmarshal gas calculator: %s", err))
		return
	}

	f.Emit(event_bus.EventLogInfo, fmt.Sprintf("SuggestSlow: %f SuggestRegular: %f SuggestPriority: %f LimitMin: %f LimitMax: %f Format: %s",
		gasCalc.SuggestSlow(),
		gasCalc.SuggestRegular(),
		gasCalc.SuggestPriority(),
		gasCalc.LimitMin(),
		gasCalc.LimitMax(),
		gasCalc.Format()),
	)

	labelCalcFee := tview.NewTextView().
		SetText(fmt.Sprintf("Total fee: %s", gasCalc.Format()))

	inputBaseFee := tview.NewInputField().
		SetFieldWidth(10).
		SetLabel("Base fee").
		SetText(fmt.Sprintf("%f", gasCalc.SuggestRegular()))

	inputPriorityFee := tview.NewInputField().
		SetFieldWidth(10).
		SetLabel("Priority fee"). // 0.1 ../ 0.25 / 3 / 4 / .. 5 (ETH)
		SetText(fmt.Sprintf("%f", gasCalc.SuggestPriority()))

	layoutForm.
		AddFormItem(inputBaseFee).
		AddFormItem(inputPriorityFee).
		AddFormItem(labelCalcFee).
		AddButton("Send", func() {

		}).
		AddButton("Update gas", func() {
			f.actionConfigureGas()
		})

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
