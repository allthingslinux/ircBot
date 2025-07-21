package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func sendMessage(session *discordgo.Session, channelID, content string) {
	if _, err := session.ChannelMessageSend(channelID, content); err != nil {
		log.Printf("Failed to send message to channel %s: %v", channelID, err)
	}
}
