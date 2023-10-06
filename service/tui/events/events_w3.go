package events

const (
	EventW3InternalGetConnections           = EventType(190)
	EventW3WalletAvailable        EventType = iota + 200
	EventW3WalletNotAvailable
	EventW3ConnAccepted
	EventW3ConnRejected
	EventW3CallGetBlockByNumber
)
