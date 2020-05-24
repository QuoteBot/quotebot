package command

import (
	"github.com/QuoteBot/quotebot/pkg/bot"
)

//AllBotCommands return the bot.Commands object containing all possible commands
func AllBotCommands() *bot.BotCommands {
	return &bot.BotCommands{
		MessageCommands:        messageCommands(),
		ReractionAddCommands:   reactionAddCommands(),
		ReactionRemoveCommands: reactionDeleteCommands(),
	}
}
