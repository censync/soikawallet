package api_web3

import mhda "github.com/censync/go-mhda"

type GetAccountsRequest struct {
	ChainKey mhda.ChainKey `json:"chain_key"`
}
