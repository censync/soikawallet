package api_web3

import "encoding/json"

// Request codes
const (
	reqCodePing = iota + 100
	reqCodeConnect
	reqCodeRequestAccounts
)

// Response codes
const (
	respCodePong uint16 = iota + 100
	respCodeConnectionAccepted
	respCodeConnectionRejected
	respCodeGetAccounts
	respCodeError      = 400
	respCodeErrorFatal = 501
)

type RPCMessageHeader struct {
	Version uint8  `json:"_v"`
	Id      string `json:"_id"`
	Type    uint16 `json:"type"`
}

type RPCMessageReq struct {
	RPCMessageHeader
	Data interface{} `json:"data"`
}

type RPCMessageResp struct {
	*RPCMessageHeader
	Data interface{} `json:"data"`
}

func (r *RPCMessageResp) toJSON() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return []byte("{}")
	}
	return data
}
