package command

import (
	"github.com/QuoteBot/quotebot/pkg/bot"
)

//AllBotCommands return the bot.Commands object containing all possible commands
func AllBotCommands() *bot.Commands {
	return &bot.Commands{
		MessageCommands:        messageCommands(),
		ReactionAddCommands:    reactionAddCommands(),
		ReactionRemoveCommands: reactionDeleteCommands(),
	}
}
