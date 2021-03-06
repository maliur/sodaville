package twitch

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/maliur/sodaville/database"
	"github.com/maliur/sodaville/irc"
	"github.com/maliur/sodaville/logger"
)

type TwitchBotConfig struct {
	OAuth            string
	BotName          string
	BotCommandPrefix string
	ChannelName      string
	Debug            bool
	DB               *database.Database
}

type TwitchBot struct {
	config *TwitchBotConfig
	client *irc.Client
}

func (t *TwitchBot) parseRawMessage(raw string) {
	var response string
	var err error

	if strings.Contains(raw, "PING") {
		t.client.WriteMessage("PONG :tmi.twitch.tv")
	}

	event := ParseIRCEvent(raw, t.config.BotCommandPrefix)

	// nothing to do here
	if len(event.Cmd) == 0 {
		return
	}

	switch event.Cmd {
	// check internal commands first
	case "dice":
		t.MentionUser(event.User, HandleDice())
		return
	case "cmd":
		response, err = HandleCmd(event, t.config.DB)
	default:
		// if it's not an internal command check the db
		fmt.Println(event)
		response, err = t.config.DB.GetCommandByName(event.Cmd)
	}

	if err != nil {
		log.Println(err)
		return
	}

	if len(response) != 0 {
		t.SendMessageToChannel(response)
	}
}

func NewTwitchBot(config *TwitchBotConfig, conn net.Conn) *TwitchBot {
	logger := logger.NewSlogger(log.New(os.Stdout, "", log.LstdFlags))
	t := &TwitchBot{config, nil}
	t.client = irc.NewClient(logger, conn, config.Debug, t.parseRawMessage)
	return t
}

func (t *TwitchBot) Connect() {
	t.client.WriteMessage(fmt.Sprintf("PASS oauth:%s", t.config.OAuth))
	t.client.WriteMessage(fmt.Sprintf("NICK %s", t.config.BotName))
	t.client.WriteMessage(fmt.Sprintf("JOIN #%s", t.config.ChannelName))

	t.SendMessageToChannel("/me booting up..")

	t.client.Run()
}

func (t *TwitchBot) Disconnect() {
	t.client.WriteMessage("/me intializing self destruction")
	t.client.WriteMessage("/me BOOM")
	t.client.WriteMessage(fmt.Sprintf("PART #%s", t.config.ChannelName))
	t.config.DB.Close()
}

func (t *TwitchBot) SendMessageToChannel(message string) {
	t.client.WriteMessage(fmt.Sprintf("PRIVMSG #%s :%s", t.config.ChannelName, message))
}

func (t *TwitchBot) WhisperUser(user, message string) {
	// TODO: Implement whisper user
	// t.client.WriteMessage(fmt.Sprintf("PRIVMSG #%s :/w %s %s", t.config.ChannelName, user, message))
}

func (t *TwitchBot) MentionUser(user, message string) {
	t.client.WriteMessage(fmt.Sprintf("PRIVMSG #%s :@%s %s", t.config.ChannelName, user, message))
}
