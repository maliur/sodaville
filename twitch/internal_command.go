package twitch

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/maliur/sodaville/database"
)

type InternalCommand struct {
	security bool
	name     string
	execute  func(message string) (string, error)
}

// sqlite> PRAGMA table_info(command);
// 0|id|INTEGER|1||1
// 1|name|VARCHAR(25)|1|''|0
// 2|security|BOOLEAN|1||0
// 3|response|TEXT|1|''|0
// sqlite> PRAGMA table_info(user);
// 0|id|INTEGER|1||1
// 1|name|VARCHAR(25)|1|''|0
// 2|trusted|BOOLEAN|1||0

func addCmd(event *IRCEvent, db *database.Database) string {
	// TODO: Add command to DB
	return fmt.Sprintf("command %s added", event.NewCmd)
}

func delCmd(event *IRCEvent, db *database.Database) string {
	// TODO: Remove command from DB
	return fmt.Sprintf("command %s deleted", event.Cmd)
}

func HandleDice(user string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("@%s %d", user, rand.Intn(100))
}

func HandleCmd(event *IRCEvent, db *database.Database) string {
	switch event.Action {
	case "add":
		return addCmd(event, db)
	case "del":
		return delCmd(event, db)
	}

	return "use this format to add or delete a command: $cmd <add|del> $command_name <text for add command if applicable>"
}
