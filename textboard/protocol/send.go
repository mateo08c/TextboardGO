package protocol

import (
	"fmt"
)

type Send struct {
	Destination string
	Value       string
}

func NewSend(destination string, value string) *Send {
	return &Send{
		Destination: destination,
		Value:       value,
	}
}

func (e *Send) ToData() ([]byte, error) {
	//SEND
	//destination:/app/map/set
	//content-length:28
	//
	//{"x":513,"y":23,"value":"d"}

	data := "SEND\n"
	data += "destination:" + e.Destination + "\n"
	data += fmt.Sprintf("content-length:%d\n", len(e.Value))
	data += "\n"
	data += e.Value
	data += "\x00"
	return []byte(data), nil
}
func (e *Send) Parse([]byte) (interface{}, error) {
	return nil, nil
}

func (e *Send) Name() PayloadName {
	return SEND
}
