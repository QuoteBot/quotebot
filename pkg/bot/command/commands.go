package command

import (
	"github.com/QuoteBot/quotebot/pkg/bot"
)

func AllBotCommands() *bot.BotCommands {
	return &bot.BotCommands{
		MessageCommands:      messageCommands(),
		ReractionAddCommands: reactionAddCommands(),
	}
}
