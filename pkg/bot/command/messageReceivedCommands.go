package command

import (
	"log"
	"strings"
	"syscall"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/bot/command/utils"
	"github.com/QuoteBot/quotebot/pkg/pagination"
	"github.com/bwmarrin/discordgo"
)

func messageCommands() map[string]bot.MessageCommand {
	return map[string]bot.MessageCommand{
		"shutdown":  shutdown,
		"ping":      ping,
		"quotebook": quotebook,
	}
}

func shutdown(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, id := range b.Conf.OwnersID {
		if m.Author.ID == id {
			b.Sc <- syscall.SIGINT
			return
		}
	}
	message := strings.Builder{}
	message.WriteString(m.Author.Mention())
	message.WriteString(" you can't do that")
	_, err := s.ChannelMessageSend(m.ChannelID, message.String())
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}
}

func ping(_ *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}
}

func quotebook(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) < 1 {
		//Todo notify user how to use it
		return
	}
	mentionedUser := m.Mentions[0]
	quotes, err := b.QuoteStore.GetQuotesFromUser(mentionedUser.ID, m.GuildID)
	if err != nil {
		log.Println("error while retrieving user quotes", err)
		return
	}

	//transform quotes in state
	state := pagination.NewState(quotes, m.Author, mentionedUser, m.ChannelID)
	//get the page
	page := state.GetCurrentPage()

	message, err := s.ChannelMessageSendEmbed(m.ChannelID, utils.EmbeddedQuotePageFactory(page, s))
	if err != nil {
		log.Println("error while sending embed", err)
		return
	}

	//register the state into the page manager
	b.PageManager.Add(message.ID, state)

	//then when it's aviable react to the message
	utils.ReactToPage(page, s, m.ChannelID, message.ID, emojiToReact())
}
