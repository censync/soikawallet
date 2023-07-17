package gas

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AlgorithmType string

type Calculator interface {
	BaseGas() float64
	SuggestSlow() float64
	SuggestRegular() float64
	SuggestPriority() float64
	LimitMin() float64
	LimitMax() uint64
	Debug() string
	LimitMaxGasFee(gasTipCap float64) float64
	FormatHumanGas(gas float64) string
	FormatHumanFiatPrice(gas float64) string
	Marshal() ([]byte, error)
}

type CalcOpts struct {
	// GasSymbol is a network defined gas unit name, e.g. "gwei", "satoshi"
	GasSymbol string `json:"gas_symbol"`
	// GasUnits is a units count per one base token, for evm = 1e9
	GasUnits uint64 `json:"gas_units"`
	// FiatSymbol is a short fiat currency suffix, e.g. $, €
	FiatSymbol string `json:"token_suffix"`
	// FiatCurrency is a fiat currency per one base token
	FiatCurrency float64 `json:"fiat_currency"`
}

/*
func (o CalcOpts) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}*/

/*
func NewGasCalculator(algorithm AlgorithmType, opts *CalcOpts) Calculator {
	switch algorithm {
	case AlgEVML1v1:
		return CalcEVML1V1{
			CalcOpts: opts,
		}
	case AlgBTCL1v1:
		return CalcBTCL1V1{
			CalcOpts: opts,
		}
	}
	return nil
}
*/

func Unmarshal(data []byte) (instance Calculator, err error) {
	var (
		tmp = struct {
			Type   AlgorithmType `json:"alg"`
			Config interface{}   `json:"config"`
		}{}
	)

	err = json.Unmarshal(data, &tmp)

	if err != nil {
		return nil, err
	}

	// TODO: Optimize serialization
	tmpCfg, _ := json.Marshal(tmp.Config)

	switch tmp.Type {
	case AlgEVML1v1:
		instance = &CalcEVML1V1{}
		err = json.Unmarshal(tmpCfg, instance)
	case AlgBTCL1v1:
		instance = &CalcBTCL1V1{}
		err = json.Unmarshal(tmpCfg, instance)
	default:
		err = errors.New(fmt.Sprintf(`undefined calculator algorithm %s`, tmp.Type))
	}

	return instance, err
}
