package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/maliur/sodaville/socket"
)

func main() {
	oauth := os.Getenv("OAUTH_TOKEN")
	botName := os.Getenv("BOT_USERNAME")
	channelName := os.Getenv("CHANNEL_NAME")

	logger := log.New(os.Stdout, "", log.LstdFlags)
	url := "ws://irc-ws.chat.twitch.tv:80"

	socket := socket.NewSocket(logger, botName, channelName, url, oauth, "$")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket.Connect()

	socket.SendMessageToChannel("Hello, is it me your looking for Kappa")

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}
}
