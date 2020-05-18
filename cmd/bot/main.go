package main

import (
	"log"
	"net"
	"os"

	"github.com/maliur/sodaville/database"
	"github.com/maliur/sodaville/twitch"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	oauth := os.Getenv("OAUTH_TOKEN")
	botName := os.Getenv("BOT_USERNAME")
	channelName := os.Getenv("CHANNEL_NAME")

	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatalf("could not connect to irc twitch: %v", err)
	}

	db, err := database.OpenDB("./sodaville.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	config := &twitch.TwitchBotConfig{
		OAuth:            oauth,
		BotName:          botName,
		BotCommandPrefix: "$",
		ChannelName:      channelName,
		Debug:            true,
		DB:               db,
	}

	t := twitch.NewTwitchBot(config, conn)
	t.Connect()
}
