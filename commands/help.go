package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/thoj/go-ircevent"
)

func init() {
	CommandMap["help"] = Help
}

// Displays available commands to the user.
func Help(session *discordgo.Session, message *discordgo.MessageCreate, irccon *irc.Connection) {
	args := strings.Fields(message.Content)
	if len(args) > 1 {
		sendMessage(session, message.ChannelID, fmt.Sprintf("Usage: `%shelp`", botPrefix))
		return
	}

	helpMessage := "**Available commands:**\n"
	for cmd := range CommandMap {
		helpMessage += fmt.Sprintf("• `%s%s`\n", botPrefix, cmd)
	}

	sendMessage(session, message.ChannelID, helpMessage)
}
