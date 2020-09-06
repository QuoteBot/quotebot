package bot

import (
	"log"
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

	if command, ok := b.Commands.MessageCommands[comm]; ok {
		//check rights
		if command.Check(b, s, m) {
			if err := command.Action(b, s, m); err != nil {
				println(err.Error())
				if _, err := s.ChannelMessageSend(m.ChannelID, command.Help()); err != nil {
					println(err.Error())
				}
				return
			}
			return
		}

		//if user do not have the good rights, send a message
		message := strings.Builder{}
		message.WriteString(m.Author.Mention())
		message.WriteString(" you can't do that")
		_, err := s.ChannelMessageSend(m.ChannelID, message.String())
		if err != nil {
			log.Println("error while sending message", err)
			return
		}
	}
}
