package protocol

type Ping struct{}

func NewPing() *Ping {
	return &Ping{}
}

func (e *Ping) ToData() ([]byte, error) {
	return []byte("\n\n"), nil
}

func (e *Ping) Parse([]byte) (interface{}, error) {
	return nil, nil
}

func (e *Ping) Name() PayloadName {
	return PING
}
