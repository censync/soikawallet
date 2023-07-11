package api_web3

import "github.com/censync/soikawallet/types"

type GetAccountsRequest struct {
	NetworkType types.NetworkType `json:"network_type"`
}
