package api_web3

import (
	"github.com/censync/soikawallet/api/dto"
)

func (c *Web3Connection) handlerConnAccepted(data interface{}) {
	d := data.(*dto.ResponseAcceptDTO)
	rpcResponse := &RPCMessageReq{
		Type: respCodeConnectionAccepted,
		Payload: map[string]interface{}{
			"instance_id": d.InstanceId,
		},
	}
	c.accepted[d.InstanceId] = true
	if conn, ok := c.hub[d.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerConnRejected(data interface{}) {
	d := data.(*dto.ResponseRejectDTO)
	rpcResponse := &RPCMessageReq{
		Type: respCodeConnectionRejected,
		Payload: map[string]interface{}{
			"instance_id": d.InstanceId,
		},
	}
	c.rejected[d.InstanceId] = true
	if conn, ok := c.hub[d.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}
