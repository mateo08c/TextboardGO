package parser

import (
	"errors"
	"fmt"
	"github.com/mateo08c/TextboardGO/textboard/protocol"
	"log"
	"strings"
)

type Callback interface {
	Handle(event *protocol.Payload)
}

type Parser struct {
	Payloads map[protocol.PayloadName]protocol.Payload
	Handlers map[protocol.PayloadName][]func(interface{}) error
}

func NewParser() *Parser {
	return &Parser{
		Payloads: protocol.NewPayloadMap(),
	}
}

func (p *Parser) OnConnected(handler func(*protocol.Connected) error) {
	p.On(protocol.CONNECTED, func(payload interface{}) error {
		if p, ok := payload.(*protocol.Connected); ok {
			return handler(p)
		}
		return errors.New("invalid payload type")
	})
}

func (p *Parser) OnError(handler func(*protocol.Error) error) {
	p.On(protocol.ERROR, func(payload interface{}) error {
		if p, ok := payload.(*protocol.Error); ok {
			return handler(p)
		}
		return errors.New("invalid payload type")
	})
}

func (p *Parser) OnPong(handler func() error) {
	p.On(protocol.PONG, func(payload interface{}) error {
		return handler()
	})
}

func (p *Parser) On(payloadName protocol.PayloadName, handler func(interface{}) error) {
	if p.Handlers == nil {
		p.Handlers = make(map[protocol.PayloadName][]func(interface{}) error)
	}

	p.Handlers[payloadName] = append(p.Handlers[payloadName], handler)

	log.Printf("Registered handler for %s\n", payloadName)
}

func (p *Parser) Emit(payloadName protocol.PayloadName, data interface{}) error {
	handlers, ok := p.Handlers[payloadName]
	if !ok {
		return nil
	}

	log.Printf("Emitting event %s\n", payloadName)
	for i, handler := range handlers {
		log.Printf("Calling handler %d\n", i)
		go func(handler func(interface{}) error) {
			err := handler(data)
			if err != nil {
				log.Printf("Error in handler: %s\n", err)
			}
		}(handler)
	}

	return nil
}

func (p *Parser) Parse(data []byte) (*protocol.PayloadName, interface{}, error) {
	if string(data) == "\n" {
		pp := protocol.PONG
		return &pp, nil, nil
	}

	payloadName := p.getPayloadName(data)

	payload, ok := p.Payloads[payloadName]
	if !ok {
		return nil, nil, fmt.Errorf("payload not found: %s", payloadName)
	}

	parsedPayload, err := payload.Parse(data)
	if err != nil {
		return nil, nil, err
	}

	return &payloadName, parsedPayload, nil
}

func (p *Parser) getPayloadName(data []byte) protocol.PayloadName {
	d := string(data)
	lines := strings.Split(d, "\n")
	firstLine := strings.Split(lines[0], " ")

	if len(firstLine) == 0 {
		return ""
	}

	return protocol.PayloadName(firstLine[0])
}
