package dto

import (
	resp "github.com/censync/soikawallet/api/responses"
	"github.com/censync/soikawallet/types"
)

type ConnectDTO struct {
	InstanceId string
	Origin     string
	RemoteAddr string
}

type RequestAccountsDTO struct {
	InstanceId  string
	Origin      string
	NetworkType types.NetworkType
}

type ResponseAcceptDTO struct {
	InstanceId string
	Chains     []*resp.ChainInfo
}

type ResponseRejectDTO struct {
	InstanceId string
	RemoteAddr string
}
