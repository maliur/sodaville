package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"

	"github.com/maliur/sodaville/database"
	"github.com/maliur/sodaville/twitch"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	oauth := os.Getenv("OAUTH_TOKEN")
	botName := os.Getenv("BOT_USERNAME")
	channelName := os.Getenv("CHANNEL_NAME")

	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("[BOOT] connecting to database")
	conn, err := sql.Open("sqlite3", "./sodaville.db")
	if err != nil {
		logger.Fatalf("[BOOT] could not connect to database: %v", err)
	}
	defer conn.Close()

	db := database.NewDatabase(logger, conn)
	chat := twitch.NewTwitch(logger, botName, channelName, oauth, db)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logger.Println("[BOOT] connecting to twitch chat")
	chat.Connect()

	for {
		<-interrupt
		log.Println("interrupt")
		chat.Close()
		return
	}
}
