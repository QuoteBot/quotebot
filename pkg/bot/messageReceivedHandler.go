package bot

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"syscall"
)

//MessageReceived Handle the recieved message
func (b *Bot) MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if strings.ToLower(m.Content) == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.ToLower(m.Content) == "shutdown" {
		//check if owner
		// send stop signal to the shutdown channel
		b.Sc <- syscall.SIGINT
	}
}
