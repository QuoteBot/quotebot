package command

import (
	"log"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/datastorage"
	"github.com/bwmarrin/discordgo"
)

func reactionAddCommands() map[string]bot.ReractionAddCommand {
	return map[string]bot.ReractionAddCommand{
		"💾": saveQuote,
	}
}

func saveQuote(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("error while geting message", err)
		return
	}

	for _, r := range message.Reactions {
		if r.Emoji.Name == "💾" {
			return
		}
	}

	timestamp, err := message.Timestamp.Parse()
	if err != nil {
		log.Println("error while parsing timestamp in saveQuote", err)
		return
	}
	quote := datastorage.Quote{
		QuoteId:   message.ID,
		GuildID:   m.GuildID,
		UserID:    message.Author.ID,
		Timestamp: timestamp,
		Content:   message.Content,
	}
	if b.QuoteStore.Save(&quote) != nil {
		log.Println("error while saving quote", err)
		return
	}
}
