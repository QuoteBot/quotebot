package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

//ReactionAdd Handle when a reaction is added to a message
func (b *Bot) ReactionAdd(s *discordgo.Session, event *discordgo.MessageReactionAdd) {

	//ignore bot reaction
	if event.UserID == s.State.User.ID {
		return
	}

	if f, ok := b.Commands.ReactionAddCommands[event.Emoji.Name]; ok {
		f(b, s, event)
	} else {
		log.Println(event.Emoji.Name, " : ", event.Emoji.ID)
	}
}
