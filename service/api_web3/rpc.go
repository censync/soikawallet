package api_web3

import "encoding/json"

const (
	respCodeConnectionPong uint8 = iota + 100
	respCodeErrorFatal
	respCodeErrorError
	respCodeConnectionAccepted
	respCodeConnectionRejected
)

type RPCMessageReq struct {
	Type    uint8                  `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type RPCMessageHeader struct {
	Version     uint8  `json:"version"`
	Type        uint8  `json:"type"`
	WalletId    string `json:"wallet_id,omitempty"`
	ExtensionId string `json:"extension_id"`
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
