package twitch

import (
	"fmt"
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

func cmdFromMessage(message, prefix string) string {
	cmdRegex := regexp.MustCompile(fmt.Sprintf(`:\%s\b[a-z]+(\b$|\s)`, prefix))
	match := cmdRegex.FindString(message)
	return strings.TrimSpace(strings.TrimPrefix(match, ":"))
}

func newCmdFromMessage(message, prefix string) string {
	cmdRegex := regexp.MustCompile(fmt.Sprintf(`:\%s[a-z]+\s([a-z]+\s[a-z]+|ls)`, prefix))
	parts := strings.Split(cmdRegex.FindString(message), " ")
	if len(parts) == 3 {
		return parts[2]
	}

	return ""
}

func userFromMessage(message string) string {
	userRegex := regexp.MustCompile(`@[\w]{3,25}`)
	return strings.TrimPrefix(userRegex.FindString(message), "@")
}

func actionFromMessage(message, prefix string) string {
	actionRegex := regexp.MustCompile(fmt.Sprintf(`:\%s[a-z]+\s([a-z]+\s[a-z]+|ls)`, prefix))
	parts := strings.Split(actionRegex.FindString(message), " ")
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

func argFromMessage(message, prefix string) string {
	actionRegex := regexp.MustCompile(fmt.Sprintf(`:\%s[a-z]+\s[a-z]+\s\w+\s.*`, prefix))
	parts := strings.Split(actionRegex.FindString(message), " ")
	if len(parts) > 3 {
		return strings.Join(parts[3:], " ")
	}

	return ""
}

func ParseIRCEvent(message, prefix string) *IRCEvent {
	return &IRCEvent{
		User:   userFromMessage(message),
		Cmd:    strings.TrimPrefix(cmdFromMessage(message, prefix), prefix),
		NewCmd: newCmdFromMessage(message, prefix),
		Action: actionFromMessage(message, prefix),
		Arg:    argFromMessage(message, prefix),
	}
}
