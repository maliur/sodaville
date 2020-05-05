package twitch

import (
	"fmt"
	"log"

	"github.com/maliur/sodaville/database"
	"github.com/maliur/sodaville/socket"
)

type Twitch struct {
	BotName     string
	ChannelName string
	OAuth       string
	logger      *log.Logger
	socket      *socket.Socket
	db          *database.Database
}

func NewTwitch(logger *log.Logger, botName, channelName, oauth string, db *database.Database) *Twitch {
	twitch := &Twitch{
		botName,
		channelName,
		oauth,
		logger,
		nil,
		db,
	}

	url := "ws://irc-ws.chat.twitch.tv:80"
	s := socket.NewSocket(
		logger,
		url,
		twitch.HandleEvent,
	)

	twitch.socket = s

	return twitch
}

func (t *Twitch) connectToChannel() {
	t.socket.SendTextMessage("[TWITCH]", fmt.Sprintf("PASS oauth:%s", t.OAuth))
	t.socket.SendTextMessage("[TWITCH]", fmt.Sprintf("NICK %s", t.BotName))
	t.socket.SendTextMessage("[TWITCH]", fmt.Sprintf("JOIN #%s", t.ChannelName))
	t.SendMessageToChannel("/me booting...")
}

func (t *Twitch) Connect() {
	t.socket.Connect()
	t.connectToChannel()
}

func (t *Twitch) Close() {
	t.SendMessageToChannel("/me shut down")
	t.socket.Close()
}

func (t *Twitch) SendMessageToChannel(message string) {
	if len(message) != 0 {
		t.socket.SendTextMessage("[TWITCH]", fmt.Sprintf("PRIVMSG #%s :%s", t.ChannelName, message))
	}
}

func (t *Twitch) HandleEvent(message string) {
	var response string
	event := ParseIRCEvent(message)
	switch event.Cmd {
	case "$cmd":
		response = HandleCmd(event, t.db)
	case "$dice":
		response = HandleDice(event.User)
	}

	if len(response) != 0 {
		t.SendMessageToChannel(response)
	}
}
