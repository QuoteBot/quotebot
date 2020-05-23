package bot

import (
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const floppy = "💾" //can use "U+1F4BE" but it should be cast in rune then in string (this is a single character)

//Bot the state of the bot
type Bot struct {
	Sc   chan os.Signal
	Conf *BotConfig
}

//MessageReceived Handle the recieved message
func (b *Bot) MessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if strings.ToLower(m.Content) == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if strings.ToLower(m.Content) == "shutdown" {
		//check if owner
		// send stop signal to the shutdown channel
		b.Sc <- syscall.SIGINT
	}
}

//GuildJoined Handle when join a guild
func (b *Bot) GuildJoined(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	s.ChannelMessageSend(event.Guild.SystemChannelID, "QuoteBot ready")
}

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

		log.Println("reacted to <", message.Author.Username, ":", timestamp.Format(time.RFC3339), message.Content, "> with floppy disck")
	default:
		log.Println("reacted to", event.MessageID, "with", event.Emoji.ID, "__", event.Emoji.Name)
	}
}
