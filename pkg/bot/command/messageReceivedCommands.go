package command

import (
	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"syscall"
)

func messageCommands() map[string]bot.MessageCommand {
	return map[string]bot.MessageCommand{
		"shutdown":  shutdown,
		"ping":      ping,
		"quotebook": quotebook,
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
	message.WriteString(" you can't do that")
	s.ChannelMessageSend(m.ChannelID, message.String())
}

func ping(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func quotebook(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	userQuotes, err := b.QuoteStore.GetQuotesFromUser(m.Author.ID, m.GuildID)
	if err != nil {
		log.Println("error while retrieving user quotes", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xfafafa, // White
		Description: "Use the reactions to navigate between pages",
		Image: &discordgo.MessageEmbedImage {
			URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail {
			URL: m.Author.AvatarURL(""),
		},
		Title:     "Quote book - the best of " + m.Author.String(),
		Footer: &discordgo.MessageEmbedFooter {
			Text:  "Requested by " + m.Author.String(),
		},
	}

	for _, q := range userQuotes.Quotes {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField {
			Name:  q.Timestamp.Format("2006-01-02"),
			Value: q.Content,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)

}
