package command

import (
	"log"

	"github.com/QuoteBot/quotebot/pkg/bot"
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
		log.Println("error while getting  message in forgetQuote: ", err)
	}

	for _, r := range message.Reactions {
		if r.Emoji.Name == "💾" {
			return
		}
	}

	if b.QuoteStore.Delete(m.MessageID, message.Author.ID, m.GuildID) != nil {
		log.Println("error while forgeting quote", err)
		return
	}
}
