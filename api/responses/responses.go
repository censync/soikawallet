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

type AddressTokenBalanceEntry struct {
	Name     string
	Symbol   string
	Contract string
	Balance  float64
}

type AddressTokensBalanceListResponse struct {
	Tokens map[uint32]*AddressTokenBalanceEntry
}
