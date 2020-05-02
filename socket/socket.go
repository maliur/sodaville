package socket

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/maliur/sodaville/command"
)

type Socket struct {
	Logger          *log.Logger
	Url             string
	RequestHeader   http.Header
	Conn            *websocket.Conn
	WebsocketDialer *websocket.Dialer
	sendMu          *sync.Mutex
	receiveMu       *sync.Mutex
	BotName         string
	ChannelName     string
	OAuth           string
	Prefix          string
}

func NewSocket(logger *log.Logger, botName, channelName, url, oauth, prefix string) *Socket {
	return &Socket{
		Logger:          logger,
		Url:             url,
		RequestHeader:   http.Header{},
		WebsocketDialer: &websocket.Dialer{},
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
		BotName:         botName,
		ChannelName:     channelName,
		OAuth:           oauth,
		Prefix:          prefix,
	}
}

func (s *Socket) Connect() {
	var err error

	s.Conn, _, err = s.WebsocketDialer.Dial(s.Url, s.RequestHeader)
	if err != nil {
		s.Logger.Fatalf("could not connect to twitch: %v\n", err)
		return
	}

	s.Logger.Println("connected to twitch!")

	s.Conn.SetPingHandler(func(appData string) error {
		return s.send(websocket.TextMessage, []byte("PONG"))
	})

	s.Conn.SetPongHandler(func(appData string) error {
		return s.send(websocket.TextMessage, []byte("PING"))
	})

	s.connectToChannel()
	command.InitCommands()

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
						s.SendTextMessage("PONG")
					}

					res, err := command.Run(message, s.Prefix)
					if err != nil {
						s.Logger.Printf("%v", err)
					}

					if len(res) > 0 {
						s.SendMessageToChannel(res)
					}
				}
			case websocket.BinaryMessage:
				s.Logger.Println("NOT IMPLEMENTED")
			}
		}
	}()
}

func (s *Socket) connectToChannel() {
	s.SendTextMessage(fmt.Sprintf("PASS oauth:%s", s.OAuth))
	s.SendTextMessage(fmt.Sprintf("NICK %s", s.BotName))
	s.SendTextMessage(fmt.Sprintf("JOIN #%s", s.ChannelName))
}

func (s *Socket) SendMessageToChannel(message string) {
	s.SendTextMessage(fmt.Sprintf("PRIVMSG #%s :%s", s.ChannelName, message))
}

func (s *Socket) SendTextMessage(message string) {
	err := s.send(websocket.TextMessage, []byte(message))
	if err != nil {
		s.Logger.Printf("could not send message to twitch: %v", err)
		return
	}

	if strings.Contains(message, "PASS") {
		// don't leak the token in the console
		s.Logger.Printf("< %s", "PASS **************")
	} else {
		s.Logger.Printf("< %s", message)
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
