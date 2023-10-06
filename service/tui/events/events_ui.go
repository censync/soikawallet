package events

const (
	EventLog EventType = iota + 100
	EventLogInfo
	EventLogSuccess
	EventLogWarning
	EventLogError

	EventUpdateCurrencies

	EventWalletInitialized
	EventWalletNoticeMessage
	EventDrawForce
	EventShowModal
	EventQuit

	// browser connector

	EventW3InternalConnections EventType = iota + 170
	EventW3Connect
	EventW3RequestAccounts
	EventW3ReqCallGetBlockByNumber

	EventW3Response
)
