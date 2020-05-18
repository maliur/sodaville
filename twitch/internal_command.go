package twitch

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/maliur/sodaville/database"
)

func lsCmd(event *IRCEvent, db *database.Database) (string, error) {
	msg, err := db.GetAllCommands()
	if err != nil {
		return "", err
	}

	return msg, nil
}

func addCmd(event *IRCEvent, db *database.Database) (string, error) {
	err := db.InsertCommand(event.NewCmd, event.Arg, false)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("command %s added", event.NewCmd), nil
}

func delCmd(event *IRCEvent, db *database.Database) (string, error) {
	err := db.DeleteCommand(event.NewCmd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("command %s deleted", event.NewCmd), nil
}

func HandleDice(user string) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("@%s %d", user, rand.Intn(100))
}

func HandleCmd(event *IRCEvent, db *database.Database) (string, error) {
	switch event.Action {
	case "ls":
		return lsCmd(event, db)
	case "add":
		return addCmd(event, db)
	case "del":
		return delCmd(event, db)
	}
	return "$cmd <action ls|add|del> <command name if add|del> <response for add command if applicable>", nil
}
