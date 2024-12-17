package protocol

type PayloadName string

const (
	CONNECT   PayloadName = "CONNECT"
	CONNECTED PayloadName = "CONNECTED"
	ERROR     PayloadName = "ERROR"
	PING      PayloadName = "PING"
	PONG      PayloadName = "PONG"
	SEND      PayloadName = "SEND"
	SUBSCRIBE PayloadName = "SUBSCRIBE"
)
