package protocol

import (
	"errors"
	"github.com/mateo08c/TextboardGO/internal/utils"
	"strings"
)

type Connected struct {
	Version   string
	HeartBeat []int
	Username  string
}

func NewConnected(username string) *Connected {
	return &Connected{
		Version:   "1.2",
		HeartBeat: []int{10000, 10000},
		Username:  username,
	}
}

func (c *Connected) ToData() ([]byte, error) {
	return nil, nil
}

func (c *Connected) Parse(data []byte) (interface{}, error) {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "version:") {
			c.Version = strings.TrimPrefix(line, "version:")
		} else if strings.HasPrefix(line, "heart-beat:") {
			heartBeatStr := strings.TrimPrefix(line, "heart-beat:")
			c.HeartBeat = utils.StringToIntSlice(heartBeatStr, ",")
		} else if strings.HasPrefix(line, "user-name:") {
			c.Username = strings.TrimPrefix(line, "user-name:")
		}
	}

	if c.Version == "" || c.Username == "" {
		return nil, errors.New("invalid connected payload")
	}

	return c, nil
}

func (c *Connected) GetHeartBeat() []int {
	return c.HeartBeat
}

func (c *Connected) GetUsername() string {
	return c.Username
}

func (c *Connected) GetVersion() string {
	return c.Version
}

func (c *Connected) Name() PayloadName {
	return CONNECTED
}
