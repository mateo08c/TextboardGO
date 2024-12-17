package protocol

type Pong struct{}

func NewPong() *Pong {
	return &Pong{}
}

func (e *Pong) ToData() ([]byte, error) {
	return nil, nil
}

func (e *Pong) Parse(data []byte) (interface{}, error) {
	return nil, nil
}

func (e *Pong) Name() PayloadName {
	return PONG
}
