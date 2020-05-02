package command

import "testing"

func TestExtractCmdFromMessage(t *testing.T) {
	messages := []struct {
		message  string
		expected string
	}{
		{":maliur!maliur@maliur.tmi.twitch.tv PRIVMSG #maliur :$dice", "$dice"},
		{":maliur!maliur@maliur.tmi.twitch.tv PRIVMSG #maliur :$hello_world", "$hello_world"},
		{":maliur!maliur@maliur.tmi.twitch.tv PRIVMSG #maliur :$helloWorld", "$helloWorld"},
		{":maliur!maliur@maliur.tmi.twitch.tv PRIVMSG #maliur :$hello-world", "$hello"},
	}

	for _, m := range messages {
		got := extractCmdFromMessage(m.message, "$")
		if got != m.expected {
			t.Errorf("extractCmdFromMessage(\"%s\") = %s; want %s", m.message, got, m.expected)
		}
	}
}
