package protocol

type Payload interface {
	ToData() ([]byte, error)
	Parse(data []byte) (interface{}, error)
	Name() PayloadName
}

func NewPayloadMap() map[PayloadName]Payload {
	return map[PayloadName]Payload{
		CONNECT:   &Connect{},
		CONNECTED: &Connected{},
		ERROR:     &Error{},
		PING:      &Ping{},
		PONG:      &Pong{},
		SEND:      &Send{},
		SUBSCRIBE: &Subscribe{},
	}
}
