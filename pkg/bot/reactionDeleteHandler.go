package bot

import (
	"github.com/bwmarrin/discordgo"
)

//ReactionDelete Handle when a reaction is deleted from a message
func (b *Bot) ReactionDelete(s *discordgo.Session, event *discordgo.MessageReactionRemove) {
	if f, ok := b.Commands.ReactionRemoveCommands[event.Emoji.Name]; ok {
		f(b, s, event)
	}
}
