package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thoj/go-ircevent"
)

var botPrefix string

// Defines the function signature for all command handlers.
type CommandHandler func(session *discordgo.Session, message *discordgo.MessageCreate, irccon *irc.Connection)

// Stores the mapping of command names to their handler functions.
var CommandMap = make(map[string]CommandHandler)

// Filters messages for the bot prefix and routes them to appropriate command handlers.
func CommandMapper(irccon *irc.Connection, prefix string) func(session *discordgo.Session, message *discordgo.MessageCreate) {
	botPrefix = prefix
	return func(session *discordgo.Session, message *discordgo.MessageCreate) {
		// Ignore bot messages and messages without the bot prefix
		if message.Author.Bot || !strings.HasPrefix(message.Content, botPrefix) {
			return
		}

		args := strings.Fields(message.Content)
		if len(args) == 0 {
			return
		}

		cmd := strings.TrimPrefix(args[0], botPrefix)
		if handler, ok := CommandMap[cmd]; ok {
			handler(session, message, irccon)
		}
	}
}
