package command

import (
	"errors"
	"log"

	"github.com/QuoteBot/quotebot/pkg/bot"
	"github.com/QuoteBot/quotebot/pkg/bot/command/utils"
	"github.com/QuoteBot/quotebot/pkg/datastorage"
	"github.com/bwmarrin/discordgo"
)

func emojiToReact() []string {
	return []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣"}
}

func fromEmojiToPos(emoji string) (int, error) {
	for pos, e := range emojiToReact() {
		if e == emoji {
			return pos, nil
		}
	}
	return -1, errors.New("not a valid emoji")
}

func reactionAddCommands() map[string]bot.ReactionAddCommand {
	res := map[string]bot.ReactionAddCommand{
		"💾":  saveQuote,
		"➡️": nextPage,
		"⬅️": prevPage,
	}
	//add selecQuote for all valid emoji to select a quote
	for _, e := range emojiToReact() {
		res[e] = selectQuote
	}
	return res
}

func saveQuote(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		log.Println("error while geting message", err)
		return
	}

	score := 0
	for _, r := range message.Reactions {

		if r.Emoji.Name == "💾" {
			score++
		}
	}

	timestamp, err := message.Timestamp.Parse()
	if err != nil {
		log.Println("error while parsing timestamp in saveQuote", err)
		return
	}
	quote := datastorage.Quote{
		QuoteId:   message.ID,
		GuildID:   m.GuildID,
		UserID:    message.Author.ID,
		Timestamp: timestamp,
		Content:   message.Content,
		Score:     score,
	}
	if b.QuoteStore.Save(&quote) != nil {
		log.Println("error while saving quote", err)
		return
	}
}

func nextPage(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	page, err := b.PageManager.GetNextPage(m.MessageID)
	if err != nil {
		//Maybe log not found for statistics
		return
	}
	embed := utils.EmbeddedQuotePageFactory(page, s)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		log.Println(err)
		return
	}
	utils.ClearAndReact(page, s, m.ChannelID, m.MessageID, emojiToReact())
}

func prevPage(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	page, err := b.PageManager.GetPreviousPage(m.MessageID)
	if err != nil {
		//Maybe log not found for statistics
		return
	}
	embed := utils.EmbeddedQuotePageFactory(page, s)
	_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, embed)
	if err != nil {
		log.Println(err)
		return
	}
	utils.ClearAndReact(page, s, m.ChannelID, m.MessageID, emojiToReact())
}

func selectQuote(b *bot.Bot, s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	page, err := b.PageManager.GetCurrentPage(m.MessageID)
	if err != nil {
		return
	}
	pos, err := fromEmojiToPos(m.Emoji.Name)
	if err != nil {
		log.Println(err)
		return
	}
	if pos >= len(page.Values) {
		log.Println("position out of range")
		return
	}

	err = utils.ReplaceByQuote(page.Values[pos], s, m.ChannelID, m.MessageID)
	if err != nil {
		log.Println(err)
		return
	}

	//when it's done delete the state from page handler
	b.PageManager.Delete(m.MessageID)
}
