package utils

import (
	"log"
	"strings"

	"github.com/QuoteBot/quotebot/pkg/datastorage"

	"github.com/QuoteBot/quotebot/pkg/pagination"
	"github.com/bwmarrin/discordgo"
)

func EmbeddedQuotePageFactory(page *pagination.Page, s *discordgo.Session) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "Quote Book of " + page.Mentioned.String(),
			IconURL: s.State.User.AvatarURL("128"),
		},
		Color:       0xfafafa, // White
		Description: "Use the reactions to navigate between pages",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: page.Mentioned.AvatarURL("64"),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Requested by " + page.Author.String(),
		},
	}

	for _, q := range page.Values {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  q.Timestamp.Format("2006-01-02"),
			Value: q.Content,
		})
	}
	return embed
}

func ReactToPage(page *pagination.Page, s *discordgo.Session, channelID string, messageID string, selections []string) {
	if page.HasPrev {
		err := s.MessageReactionAdd(channelID, messageID, "⬅️")
		if err != nil {
			log.Println("error while reacting", err)
			return
		}
	}

	if page.HasNext {
		err := s.MessageReactionAdd(channelID, messageID, "➡️")
		if err != nil {
			log.Println("error while reacting", err)
			return
		}
	}

	for i := 0; i < len(page.Values); i++ {
		err := s.MessageReactionAdd(channelID, messageID, selections[i])
		if err != nil {
			log.Println("error while reacting", err)
			return
		}
	}
}

func ClearAndReact(page *pagination.Page, s *discordgo.Session, channelID string, messageID string, selections []string) {
	s.MessageReactionsRemoveAll(channelID, messageID)
	ReactToPage(page, s, channelID, messageID, selections)
}

func ReplaceByQuote(quote datastorage.Quote, s *discordgo.Session, channelID string, messageID string) error {
	quoteAuthor, err := s.User(quote.UserID)
	if err != nil {
		return err
	}

	builder := strings.Builder{}
	builder.WriteString("> ")
	builder.WriteString(quote.Content)
	builder.WriteRune('\n')
	builder.WriteString(quoteAuthor.Mention())

	//new message with the quote
	_, err = s.ChannelMessageSend(channelID, builder.String())
	if err != nil {
		return err
	}

	//remove Embedded
	err = s.ChannelMessageDelete(channelID, messageID)
	if err != nil {
		return err
	}
	return nil
}
