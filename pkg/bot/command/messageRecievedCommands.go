package command

import (
	"strings"
	"syscall"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/bwmarrin/discordgo"
)

func MessageCommands() map[string]bot.MessageCommand {
	return map[string]bot.MessageCommand{
		"shutdown": shutdown,
		"ping":     ping,
	}
}

func shutdown(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, id := range b.Conf.OwnersID {
		if m.Author.ID == id {
			b.Sc <- syscall.SIGINT
			return
		}
	}
	message := strings.Builder{}
	message.WriteString(m.Author.Mention())
	message.WriteString(" you are not allowed to shutdown me")
	s.ChannelMessageSend(m.ChannelID, message.String())
}

func ping(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
