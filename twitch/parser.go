package twitch

import (
	"regexp"
	"strings"
)

type IRCEvent struct {
	User   string
	Cmd    string
	NewCmd string
	Action string
	Arg    string
}

func cmdFromMessage(message string) string {
	cmdRegex := regexp.MustCompile(`:\$\b[a-z]+(\b$|\s)`)
	match := cmdRegex.FindString(message)
	return strings.TrimSpace(strings.TrimPrefix(match, ":"))
}

func newCmdFromMessage(message string) string {
	cmdRegex := regexp.MustCompile(`[^:]\$\b[a-z]+(\b$|\s)`)
	match := cmdRegex.FindString(message)
	return strings.TrimSpace(match)
}

func userFromMessage(message string) string {
	userRegex := regexp.MustCompile(`@[\w]{3,25}`)
	return strings.TrimPrefix(userRegex.FindString(message), "@")
}

func actionFromMessage(message string) string {
	actionRegex := regexp.MustCompile(`:\$[a-z]+\s[a-z]+\s\$[a-z]+`)
	parts := strings.Split(actionRegex.FindString(message), " ")
	if len(parts) >= 3 {
		return parts[1]
	}

	return ""
}

func argFromMessage(message string) string {
	actionRegex := regexp.MustCompile(`:\$[a-z]+\s[a-z]+\s\$\w+\s.*`)
	parts := strings.Split(actionRegex.FindString(message), " ")
	if len(parts) > 3 {
		return strings.Join(parts[3:], " ")
	}

	return ""
}

func ParseIRCEvent(message string) *IRCEvent {
	return &IRCEvent{
		User:   userFromMessage(message),
		Cmd:    cmdFromMessage(message),
		NewCmd: newCmdFromMessage(message),
		Action: actionFromMessage(message),
		Arg:    argFromMessage(message),
	}
}
