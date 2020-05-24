package command

import (
	"log"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/datastorage"
	"github.com/bwmarrin/discordgo"
)

func reactionDeleteCommands() map[string]bot.ReractionRemoveCommand {
	return map[string]bot.ReractionRemoveCommand{
		"💾": forgetQuote,
	}
}

func forgetQuote(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("error while geting message in forgetQuote", err)
		return
	}
	timestamp, err := message.Timestamp.Parse()
	if err != nil {
		log.Println("error while parsing timestamp in forgetQuote", err)
		return
	}
	quote := &datastorage.Quote{
		QuoteId:   message.ID,
		GuildID:   m.GuildID,
		UserID:    message.Author.ID,
		Timestamp: timestamp,
		Content:   message.Content,
	}
	if b.QuoteStore.Forget(quote) != nil {
		log.Println("error while forgeting quote", err)
		return
	}
}
