package bot

import (
	"github.com/bwmarrin/discordgo"
)

//ReactionAdd Handle when a reaction is added to a message
func (b *Bot) ReactionAdd(s *discordgo.Session, event *discordgo.MessageReactionAdd) {
	if f, ok := b.Commands.ReractionAddCommands[event.Emoji.Name]; ok {
		f(b, s, event)
	}
}
