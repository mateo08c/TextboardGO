package protocol

import "fmt"

type Subscribe struct {
	Id          string
	Destination string
}

func NewSubscribe(id int, destination string) *Subscribe {
	return &Subscribe{
		Id:          fmt.Sprintf("sub-%d", id),
		Destination: destination,
	}
}

func (e *Subscribe) ToData() ([]byte, error) {
	//SUBSCRIBE
	//id:sub-1
	//destination:/app/map/chunk/0/-2

	data := "SUBSCRIBE\n"
	data += "id:" + e.Id + "\n"
	data += "destination:" + e.Destination + "\n"
	data += "\n"
	data += "\x00"

	fmt.Printf("%s\n", data)

	return nil, nil
}

func (e *Subscribe) Parse(data []byte) (interface{}, error) {
	return nil, nil
}

func (e *Subscribe) Name() PayloadName {
	return SUBSCRIBE
}
