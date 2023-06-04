package api_web3

type ResponsePong struct {
	WalletStatus uint8 `json:"wallet_status"`
}

type ResponseConnectionAccepted struct {
	InstanceId string `json:"instance_id"`
}

type ResponseErrorFatal struct {
	Error string `json:"error"`
}
