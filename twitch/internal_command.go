package twitch

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/maliur/sodaville/database"
)

// type InternalCommand struct {
// 	security bool
// 	name     string
// 	execute  func(message string) (string, error)
// }

func addCmd(event *IRCEvent, db *database.Database) string {
	// TODO: Add command to DB
	return fmt.Sprintf("command %s added", event.Cmd)
}

func delCmd(event *IRCEvent, db *database.Database) string {
	// TODO: Remove command from DB
	return fmt.Sprintf("command %s deleted", event.Cmd)
}

func HandleDice() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(100))
}

func HandleCmd(event *IRCEvent, db *database.Database) string {
	switch event.Action {
	case "add":
		return addCmd(event, db)
	case "del":
		return delCmd(event, db)
	}

	return "no action found"
}
