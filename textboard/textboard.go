package textboard

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mateo08c/TextboardGO/textboard/parser"
	"github.com/mateo08c/TextboardGO/textboard/protocol"
	"log"
	"sync"
	"time"
)

type Textboard struct {
	options  *Options
	parser   *parser.Parser
	conn     *websocket.Conn
	writeMux sync.Mutex
	wg       sync.WaitGroup
	stop     chan struct{}
	status   Status
}

type Options struct {
	Authorization     string
	WebsocketSecure   bool
	WebsocketHostname string
	WebsocketPort     int
	PingPeriod        int
	WriteWait         int
	ReadTimeout       int
	WebsocketPath     string
}

// build websocket connection
func (t *Textboard) websocketURL() string {
	p := "ws"
	if t.options.WebsocketSecure {
		p = "wss"
	}

	return fmt.Sprintf("%s://%s:%d%s", p, t.options.WebsocketHostname, t.options.WebsocketPort, t.options.WebsocketPath)
}

func (t *Textboard) GetWebsocketURL() string {
	return t.websocketURL()
}

func NewTextboard(o *Options) *Textboard {
	t := &Textboard{
		options:  o,
		writeMux: sync.Mutex{},
		parser:   parser.NewParser(),
		wg:       sync.WaitGroup{},
		stop:     make(chan struct{}),
		status:   Disconnected,
	}

	t.parser.OnConnected(func(c *protocol.Connected) error {
		log.Printf("Connected to Textboard v%s as %s\n", c.GetVersion(), c.GetUsername())

		t.status = Authenticated

		return nil
	})

	t.parser.OnError(func(e *protocol.Error) error {
		log.Printf("Error: %s\n", e.Message)
		return nil
	})

	t.parser.OnPong(func() error {
		return t.Ping()
	})

	return t
}

func (t *Textboard) WriteLetter(x, y int, letter string) error {
	type Write struct {
		X     int    `json:"x"`
		Y     int    `json:"y"`
		Value string `json:"value"`
	}

	e, err := json.Marshal(Write{
		X:     x,
		Y:     y,
		Value: letter,
	})

	log.Printf("Sending letter to %d, %d: %s\n", x, y, letter)

	if err != nil {
		return err
	}

	s := protocol.NewSend("/app/map/set", string(e))
	err = t.Write(s)
	if err != nil {
		log.Printf("Error sending letter: %s\n", err)
		return err
	}

	return nil
}

func (t *Textboard) Ping() error {
	p := protocol.NewPing()
	err := t.Write(p)
	if err != nil {
		log.Printf("Error sending ping: %s\n", err)
	}

	return err
}

func (t *Textboard) Connect() error {
	log.Printf("Connecting to %s\n", t.websocketURL())
	conn, _, err := websocket.DefaultDialer.Dial(t.websocketURL(), nil)
	if err != nil {
		return err
	}

	t.conn = conn

	log.Printf("Connected to %s\n", t.websocketURL())

	c := protocol.NewConnect(t.options.Authorization)
	t.status = Authenticating
	err = t.Write(c)
	if err != nil {
		return err
	}

	err = t.Read()
	if err != nil {
		log.Printf("Error reading: %s\n", err)
	}

	return nil
}

func (t *Textboard) startPingTimer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-t.stop:
			log.Printf("Stopping ping timer\n")
			return
		case <-ticker.C:
			log.Printf("Sending ping\n")
			err := t.Write(protocol.NewPing())
			if err != nil {
				log.Printf("Error sending ping: %s\n", err)
			}
		}
	}
}

func (t *Textboard) Close() {
	log.Printf("Closing connection\n")
	close(t.stop)

	t.wg.Wait()
	log.Printf("Connection closed successfully\n")
}

func (t *Textboard) Disconnect() error {
	return nil
}

func (t *Textboard) Write(payload protocol.Payload) error {
	defer t.writeMux.Unlock()
	t.writeMux.Lock()

	data, err := payload.ToData()
	if err != nil {
		return err
	}

	log.Printf("Sending payload %s\n", payload.Name())

	err = t.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}

func (t *Textboard) Read() error {
	t.wg.Add(1)

	defer func() {
		t.wg.Done()
	}()

	log.Printf("Starting read loop\n")

	for {
		select {
		default:
			_, data, err := t.conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Connection closed gracefully: %v\n", err)
					return nil
				}

				log.Printf("Error reading message: %v\n", err)
				return err
			}

			var pn *protocol.PayloadName
			var p interface{}

			pn, p, err = t.parser.Parse(data)
			if err != nil {
				log.Printf("Error parsing payload: %s\n", err)
				continue
			}

			log.Printf("Received payload %s\n", *pn)

			err = t.parser.Emit(*pn, p)
			if err != nil {
				log.Printf("Error emitting payload: %s\n", err)
				continue
			}

		case <-t.stop:
			log.Printf("Stopping read loop\n")
			return nil
		}
	}
}

func (t *Textboard) Handle() error {
	return nil
}

func (t *Textboard) GetStatus() Status {
	return t.status
}
