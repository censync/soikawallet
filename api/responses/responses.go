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
