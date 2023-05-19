package dto

type ConnectDTO struct {
	InstanceId string
	Origin     string
	RemoteAddr string
}

type RequestAccountsDTO struct {
	InstanceId string
	Origin     string
	Network    string
}

type ResponseAcceptDTO struct {
	InstanceId string
	RemoteAddr string
}

type ResponseRejectDTO struct {
	InstanceId string
	RemoteAddr string
}
