package responses

import (
	mhda "github.com/censync/go-mhda"
)

// exported type
type AddressResponse struct {
	Address        string
	Path           string
	IsExternal     bool
	ChainKey       mhda.ChainKey
	AddressIndex   mhda.AddressIndex
	Account        mhda.AccountIndex
	Label          string
	IsW3           bool
	IsKeyDelivered bool
}

type AccountResponse struct {
	// Path        string
	ChainKey mhda.ChainKey
	Account  mhda.AccountIndex
	Label    string
}

type AirGapMessage struct {
	Chunks []string
}

type AddressTokenEntry struct {
	Standard uint8
	Name     string
	Symbol   string
	Contract string
}

type BaseCurrency struct {
	Symbol   string
	Decimals int
}

type AddressTokensListResponse map[string]*AddressTokenEntry

type AddressTokenBalanceEntry struct {
	Token   *AddressTokenEntry
	Balance float64
}

type AddressTokensBalanceListResponse map[string]*AddressTokenBalanceEntry

type TokenConfig struct {
	Standard uint8
	Name     string
	Symbol   string
	Contract string
	Decimals int
}

type ChainInfo struct {
	Name     string        `json:"name"`
	ChainKey mhda.ChainKey `json:"chain_key"`
}

type CalculatorConfig struct {
	Calculator []byte
}
