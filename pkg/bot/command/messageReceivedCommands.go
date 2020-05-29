package command

import (
	"log"
	"strings"
	"syscall"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/pagination"
	"github.com/bwmarrin/discordgo"
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
	_, err := s.ChannelMessageSend(m.ChannelID, message.String())
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}
}

func ping(_ *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}
}

func quotebook(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	mentionedUser := m.Mentions[0]
	quotes, err := b.QuoteStore.GetQuotesFromUser(mentionedUser.ID, m.GuildID)
	if err != nil {
		log.Println("error while retrieving user quotes", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Quote Book of " + mentionedUser.String(),
			IconURL: s.State.User.AvatarURL("128"),
		},
		Color:       0xfafafa, // White
		Description: "Use the reactions to navigate between pages",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: mentionedUser.AvatarURL("64"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Requested by " + m.Author.String(),
		},
	}

	//transform quotes in state
	state := pagination.NewState(quotes)
	//get the page
	page := state.GetCurrentPage()

	for _, q := range page.Values {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  q.Timestamp.Format("2006-01-02"),
			Value: q.Content,
		})
	}

	message, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}

	//register the state into the page manager
	b.PageHandler.Add(message.ID, state)

	if page.HasPrev {
		err = s.MessageReactionAdd(m.ChannelID, message.ID, "⬅️")
		if err != nil {
			log.Println("error while sending embed", err)
			return
		}
	}

	if page.HasNext {
		err = s.MessageReactionAdd(m.ChannelID, message.ID, "➡️")
		if err != nil {
			log.Println("error while sending embed", err)
			return
		}
	}
}
