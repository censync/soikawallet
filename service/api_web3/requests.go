package api_web3

type RequestConnect struct {
	InstanceId string `json:"instance_id"`
	Origin     string `json:"origin"`
	RemoteAddr string `json:"remote_addr"`
}

type RequestRequestAccounts struct {
	InstanceId string `json:"instance_id"`
	Origin     string `json:"origin"`
	Network    string `json:"network"`
}
