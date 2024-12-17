package protocol

import (
	"github.com/mateo08c/TextboardGO/internal/utils"
	"strings"
)

type Connect struct {
	Authorization string
	AcceptVersion []string
	HeartBeat     []int
}

func NewConnect(authorization string) Payload {
	return &Connect{
		Authorization: authorization,
		AcceptVersion: []string{"1.2", "1.1", "1.0"},
		HeartBeat:     []int{10000, 10000},
	}
}

func (c *Connect) ToData() ([]byte, error) {
	//CONNECT
	//Authorization: xxx
	//accept-version:1.2,1.1,1.0
	//heart-beat:10000,10000

	data := "CONNECT\n"
	data += "Authorization:" + c.Authorization + "\n"
	data += "accept-version:" + strings.Join(c.AcceptVersion, ",") + "\n"
	data += "heart-beat:" + utils.IntSliceToString(c.HeartBeat, ",") + "\n"
	data += "\n"
	data += "\x00"

	return []byte(data), nil
}

func (c *Connect) Parse(data []byte) (interface{}, error) {
	return nil, nil
}

func (c *Connect) Name() PayloadName {
	return CONNECT
}
