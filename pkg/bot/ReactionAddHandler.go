package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

//ReactionAdd Handle when a reaction is added to a message
func (b *Bot) ReactionAdd(s *discordgo.Session, event *discordgo.MessageReactionAdd) {

	switch event.Emoji.Name {
	case floppy:
		// Refactor
		message, err := s.ChannelMessage(event.ChannelID, event.MessageID)
		if err != nil {
			log.Println(err)
			return
		}
		//message.EditedTimestamp is a decorated string
		log.Println("this is the timestamp", message.Timestamp)
		timestamp, err := message.Timestamp.Parse()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("reacted to <", message.Author.Username, ":", timestamp.Format(time.RFC3339), message.Content, "> with floppy disk")
	default:
		log.Println("reacted to", event.MessageID, "with", event.Emoji.ID, "__", event.Emoji.Name)
	}
}
