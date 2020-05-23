package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

//MessageReceived Handle the recieved message
func (b *Bot) MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	comm := strings.ToLower(strings.Split(m.Content, " ")[0])

	for c, f := range b.Commands.MessageCommands {
		if c == comm {
			f(b, s, m)
			return
		}
	}
}
