package api_web3

import (
	"encoding/json"
	"github.com/censync/soikawallet/api/dto"
)

func (c *Web3Connection) handlerWalletAvailable(data interface{}) {
	c.walletId = data.(string)
	rpcResponse := c.newRPCResponse(respCodePong, map[string]interface{}{
		"wallet_status": c.walletStatus(),
	})
	for _, conn := range c.hub {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerWalletPing(data interface{}) {
	instanceId := data.(string)
	rpcResponse := c.newRPCResponse(respCodePong, map[string]interface{}{
		"wallet_status": c.walletStatus(),
	})
	if conn, ok := c.hub[instanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerConnAccepted(data interface{}) {
	d := data.(*dto.ResponseAcceptDTO)
	rpcResponse := c.newRPCResponse(respCodeConnectionAccepted, map[string]interface{}{
		"wallet_status": c.walletStatus(),
		"chains":        d.Chains,
	})
	c.accepted[d.InstanceId] = true
	if conn, ok := c.hub[d.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerConnRejected(data interface{}) {
	d := data.(*dto.ResponseRejectDTO)
	rpcResponse := c.newRPCResponse(respCodeConnectionRejected, map[string]interface{}{})
	c.rejected[d.InstanceId] = true
	if conn, ok := c.hub[d.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerAccountsGet(data interface{}) {
	d := data.(*dto.RequestAccountsDTO)

	rpcResponse := c.newRPCResponse(respCodeGetAccounts, map[string]interface{}{
		"wallet_status": c.walletStatus(),
	})
	if conn, ok := c.hub[d.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}

func (c *Web3Connection) handlerCallGetBlockByNumber(data interface{}) {
	p := data.(*dto.ResponseGetBlockByNumberDTO)
	m := map[string]interface{}{}
	json.Unmarshal(p.Data, &m)
	rpcResponse := c.newRPCResponse(respCodeGetBlockByNumber, map[string]interface{}{
		"block": m,
	})
	if conn, ok := c.hub[p.InstanceId]; ok {
		_ = conn.WriteJSON(rpcResponse)
	}
}
