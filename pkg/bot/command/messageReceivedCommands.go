package command

import (
	"errors"
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
		"shutdown":  &shutdown{},
		"ping":      &ping{},
		"quotebook": &quotebook{},
		"help":      &help{},
	}
}

//Shutdown section
type shutdown struct{}

func (sh *shutdown) Action(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) error {
	b.Sc <- syscall.SIGINT
	return nil
}

func (sh *shutdown) Check(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	for _, id := range b.Conf.OwnersID {
		if m.Author.ID == id {
			return true
		}
	}
	return false
}

func (sh *shutdown) Help() string {
	return "shutdown : if you are a owner of this bot, you can shutdown it"
}

//Ping section

type ping struct{}

func (ping *ping) Action(_ *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
	if err != nil {
		log.Println("error while sending embed", err)
		return err
	}
	return nil
}

func (ping *ping) Check(_ *bot.Bot, _ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return true
}

func (ping *ping) Help() string {
	return "ping : ping the bot, it should respond with pong"
}

//quotebook section
type quotebook struct{}

func (qu *quotebook) Action(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if len(m.Mentions) < 1 {
		return errors.New("not enougth arguments")
	}
	mentionedUser := m.Mentions[0]
	quotes, err := b.QuoteStore.GetQuotesFromUser(mentionedUser.ID, m.GuildID)
	if err != nil {
		log.Println("error while retrieving user quotes", err)
		return err
	}

	//transform quotes in state
	state := pagination.NewState(quotes, m.Author, mentionedUser, m.ChannelID)
	//get the page
	page := state.GetCurrentPage()

	message, err := s.ChannelMessageSendEmbed(m.ChannelID, utils.EmbeddedQuotePageFactory(page, s))
	if err != nil {
		log.Println("error while sending embed", err)
		return err
	}

	//register the state into the page manager
	b.PageManager.Add(message.ID, state)

	//then when it's aviable react to the message
	utils.ReactToPage(page, s, m.ChannelID, message.ID, emojiToReact())
	return nil
}

func (qu *quotebook) Check(_ *bot.Bot, _ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return true
}

func (qu *quotebook) Help() string {
	return "quotebook @user : open the quotebook of the notified user"
}

//help section
type help struct{}

func (h *help) Action(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageCreate) error {
	message := strings.Builder{}
	for _, v := range messageCommands() {
		message.WriteString(v.Help())
		message.WriteRune('\n')
		message.WriteRune('\n')
	}

	if _, err := s.ChannelMessageSend(m.ChannelID, message.String()); err != nil {
		return err
	}
	return nil
}

func (h *help) Check(_ *bot.Bot, _ *discordgo.Session, _ *discordgo.MessageCreate) bool {
	return true
}

func (h *help) Help() string {
	return "help [command]: show help and informations"
}
