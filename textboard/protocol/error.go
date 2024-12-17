package protocol

import (
	"bytes"
	"strconv"
	"strings"
)

type Error struct {
	Message       string
	ContentLength int
}

func NewError() *Error {
	return &Error{}
}

func (e *Error) ToData() ([]byte, error) {
	return nil, nil
}

func (e *Error) Parse(data []byte) (interface{}, error) {
	// Split the data into lines
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(string(parts[0]))
		value := strings.TrimSpace(string(parts[1]))

		switch key {
		case "message":
			e.Message = value
		case "content-length":
			length, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			e.ContentLength = length
		}
	}
	return e, nil
}

func (e *Error) Name() PayloadName {
	return ERROR
}
