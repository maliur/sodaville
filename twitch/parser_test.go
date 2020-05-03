package twitch

import "testing"

func TestUserFromMessage(t *testing.T) {
	messages := []struct {
		message  string
		expected string
	}{
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :hello", "john"},
		{":john_doe!john_doe@john_doe.tmi.twitch.tv PRIVMSG #sodaville :hello", "john_doe"},
	}

	for _, m := range messages {
		got := userFromMessage(m.message)
		if got != m.expected {
			t.Errorf("userFromMessage(\"%s\") = %s; want %s", m.message, got, m.expected)
		}
	}
}

func TestCmdFromMessage(t *testing.T) {
	messages := []struct {
		message  string
		expected string
	}{
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$dice", "$dice"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$hello_world", ""},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$helloWorld", ""},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$HeLlo", ""},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$hello-world", ""},
	}

	for _, m := range messages {
		got := cmdFromMessage(m.message)
		if got != m.expected {
			t.Errorf("cmdFromMessage(\"%s\") = %s; want %s", m.message, got, m.expected)
		}
	}
}

func TestActionFromMessage(t *testing.T) {
	messages := []struct {
		message  string
		expected string
	}{
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd add $yolo", "add"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd del $yolo", "del"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd add $yolo You only live once", "add"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd ADD $yolo", ""},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd add", ""},
	}

	for _, m := range messages {
		got := actionFromMessage(m.message)
		if got != m.expected {
			t.Errorf("actionFromMessage(\"%s\") = %s; want %s", m.message, got, m.expected)
		}
	}
}

func TestArgFromMessage(t *testing.T) {
	messages := []struct {
		message  string
		expected string
	}{
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd add $yolo You only live once", "You only live once"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd del $yolo Hello, World", "Hello, World"},
		{":john!john@john.tmi.twitch.tv PRIVMSG #sodaville :$cmd ADD $yolo Hello, World", ""},
	}

	for _, m := range messages {
		got := argFromMessage(m.message)
		if got != m.expected {
			t.Errorf("argFromMessage(\"%s\") = %s; want %s", m.message, got, m.expected)
		}
	}
}
