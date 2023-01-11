package responses

import "github.com/censync/soikawallet/types"

// exported type
type AddressResponse struct {
	Address           string
	Path              string
	IsExternal        bool
	AddressIndex      types.AddressIndex
	IsHardenedAddress bool
	CoinType          types.CoinType
	Account           types.AccountIndex
}

type AirGapMessageResponse struct {
	Chunks []string
}

type AddressTokenEntry struct {
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
