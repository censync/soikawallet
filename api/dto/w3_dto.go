package dto

import (
	mhda "github.com/censync/go-mhda"
	resp "github.com/censync/soikawallet/api/responses"
)

type ConnectDTO struct {
	InstanceId string
	Origin     string
	RemoteAddr string
}

type RequestAccountsDTO struct {
	InstanceId string
	Origin     string
	ChainKey   mhda.ChainKey
}

type ResponseAcceptDTO struct {
	InstanceId string
	Chains     []*resp.ChainInfo
}

type ResponseRejectDTO struct {
	InstanceId string
	RemoteAddr string
}
