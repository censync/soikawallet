package api_web3

import mhda "github.com/censync/go-mhda"

type GetAccountsRequest struct {
	ChainKey mhda.ChainKey `json:"chain_key"`
}

type RPCRequest struct {
	ChainKey mhda.ChainKey `json:"chain_key"`
	Method   string        `json:"request"`
	Params   []string      `json:"params"`
}
