package socket

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Socket struct {
	Logger          *log.Logger
	Url             string
	RequestHeader   http.Header
	Conn            *websocket.Conn
	WebsocketDialer *websocket.Dialer
	IncomingHandler func(message string)
	sendMu          *sync.Mutex
	receiveMu       *sync.Mutex
}

func NewSocket(logger *log.Logger, url string, incomingHandler func(message string)) *Socket {
	return &Socket{
		Logger:          logger,
		Url:             url,
		RequestHeader:   http.Header{},
		WebsocketDialer: &websocket.Dialer{},
		IncomingHandler: incomingHandler,
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
	}
}

func (s *Socket) Connect() {
	var err error

	s.Conn, _, err = s.WebsocketDialer.Dial(s.Url, s.RequestHeader)
	if err != nil {
		s.Logger.Fatalf("could not connect to twitch: %v\n", err)
		return
	}

	s.Logger.Println("Connecting...")

	s.Conn.SetPingHandler(func(appData string) error {
		return s.send(websocket.TextMessage, []byte("PONG"))
	})

	s.Conn.SetPongHandler(func(appData string) error {
		return s.send(websocket.TextMessage, []byte("PING"))
	})

	go func() {
		for {
			s.receiveMu.Lock()
			messageType, buffer, err := s.Conn.ReadMessage()
			s.receiveMu.Unlock()
			if err != nil {
				s.Logger.Println("read:", err)
				continue
			}

			switch messageType {
			case websocket.TextMessage:
				{
					message := string(buffer)
					s.Logger.Printf("> %s", message)

					if strings.Contains(message, "PING") {
						s.SendTextMessage("[SOCKET]", "PONG")
						continue
					}

					s.IncomingHandler(message)
				}
			case websocket.BinaryMessage:
				s.Logger.Println("NOT IMPLEMENTED")
			}
		}
	}()
}

func (s *Socket) SendTextMessage(prefix, message string) {
	err := s.send(websocket.TextMessage, []byte(message))
	if err != nil {
		s.Logger.Printf("could not send message to twitch: %v", err)
		return
	}

	if strings.Contains(message, "PASS") {
		// don't leak the token in the console
		s.Logger.Printf("< %s", "PASS **************")
	} else {
		if len(prefix) > 0 {
			prefix = prefix + " "
		}

		s.Logger.Printf("%s< %s", prefix, message)
	}
}

func (s *Socket) send(messageType int, data []byte) error {
	s.sendMu.Lock()
	err := s.Conn.WriteMessage(messageType, data)
	s.sendMu.Unlock()

	return err
}

func (s *Socket) Close() {
	err := s.send(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		s.Logger.Println("write close:", err)
	}
	s.Conn.Close()
}
