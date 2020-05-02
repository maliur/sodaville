package command

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Move to DB
var commands []command
var trusted = []string{
	"maliur",
}

type cmdType int

const (
	text cmdType = iota
	code
)

type command struct {
	security bool
	name     string
	typ      cmdType
	response string
	execute  func(message string) (string, error)
}

func extractCmdFromMessage(message, prefix string) string {
	cmdRegex := regexp.MustCompile(`:\` + prefix + `[A-z]+`)
	return strings.TrimPrefix(cmdRegex.FindString(message), ":"+prefix)
}

func extractUsernameFromMessage(message string) string {
	usernameRegex := regexp.MustCompile(`@[a-z]+`)
	return strings.TrimPrefix(usernameRegex.FindString(message), "@")
}

func isTrusted(username string) bool {
	for _, t := range trusted {
		if t == username {
			return true
		}
	}

	return false
}

func InitCommands() {
	commands = append(commands, command{false, "dice", code, "", dice})
	commands = append(commands, command{true, "cmd", code, "", handleCmd})
}

func Run(message, prefix string) (string, error) {
	username := extractUsernameFromMessage(message)
	cmdName := extractCmdFromMessage(message, prefix)
	for _, cmd := range commands {
		if cmd.name == cmdName {
			switch cmd.typ {
			case code:
				{
					if cmd.security && !isTrusted(username) {
						// not trusted to execute command
						return "", nil
					}

					msg, err := cmd.execute(message)
					if err != nil {
						return "", err
					}

					return msg, nil
				}
			case text:
				return cmd.response, nil
			}
		}
	}

	if len(cmdName) > 0 {
		return "", fmt.Errorf("[WARN] no command found with name: %s", cmdName)
	}

	return "", nil
}

func dice(_ string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(100)), nil
}

func isValidCommand(cmd string) bool {
	cmdRegex := regexp.MustCompile(`\$[A-z]+`)
	return cmdRegex.MatchString(cmd)
}

func addCmd(parts []string) (string, error) {
	wrongFormat := "the format for adding a command is: $cmd add $command_name <text>"
	if len(parts) < 4 {
		return wrongFormat, nil
	}

	cmdName := parts[2]
	if !isValidCommand(cmdName) {
		return wrongFormat, nil
	}

	cmdResponse := strings.Join(parts[3:], " ")
	commands = append(commands, command{false, strings.TrimPrefix(cmdName, "$"), text, cmdResponse, nil})

	return fmt.Sprintf("command %s added", cmdName), nil
}

func delCmd(parts []string) (string, error) {
	wrongFormat := "the format for deleting a command is: $cmd del $command_name"
	if len(parts) != 3 {
		return wrongFormat, nil
	}

	cmdName := parts[2]
	if !isValidCommand(cmdName) {
		return wrongFormat, nil
	}

	index := -1
	for idx, cmd := range commands {
		if cmd.name == strings.TrimPrefix(cmdName, "$") {
			index = idx
			break
		}
	}

	if index == -1 {
		return fmt.Sprintf("command %s was not found", cmdName), nil
	}

	commands[index] = commands[len(commands)-1]
	commands = commands[:len(commands)-1]

	return fmt.Sprintf("command %s deleted", cmdName), nil
}

func handleCmd(message string) (string, error) {
	addRegex := regexp.MustCompile(`:\$[A-z]+\sadd\s\$[A-z]+\s.*`)
	delRegex := regexp.MustCompile(`:\$[A-z]+\sdel\s\$[A-z]+`)
	var match string
	if addRegex.MatchString(message) {
		match = addRegex.FindString(message)
	} else if delRegex.MatchString(message) {
		match = delRegex.FindString(message)
	} else {
		match = ""
	}

	parts := strings.Split(match, " ")

	action := parts[1]
	switch action {
	case "add":
		return addCmd(parts)
	case "del":
		return delCmd(parts)
	}

	return "no action found", nil
}
