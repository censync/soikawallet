package api_web3

import "github.com/censync/soikawallet/types"

type GetAccountsRequest struct {
	CoinType types.CoinType `json:"coin_type"`
}
