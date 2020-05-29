package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

//MessageReceived Handle the recieved message
func (b *Bot) MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	prefix := b.Conf.Prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, prefix) {
		return
	}

	comm := m.Content[len(prefix):]
	comm = strings.ToLower(strings.Split(comm, " ")[0])

	for c, f := range b.Commands.MessageCommands {
		if c == comm {
			f(b, s, m)
			return
		}
	}
}
