package api_web3

type ResponsePong struct {
	WalletState uint8 `json:"wallet_state"`
}

type ResponseConnectionAccepted struct {
	InstanceId string `json:"instance_id"`
}

type ResponseErrorFatal struct {
	Error string `json:"error"`
}
