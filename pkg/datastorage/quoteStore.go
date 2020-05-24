package datastorage

import (
	"time"
)

type Quote struct {
	UserID    string    `json:"userID"`
	GuildID   string    `json:"guildId"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

type QuoteStore interface {
	Save(quote *Quote) error
	Forget(quote *Quote) error
	//GetQuotesFromUser(userID string, guildID string) []Quote
	//GetAllQuotes(guildID string) []Quote
	//FindQuotesFromUser(userID string, guildID string, search string) []Quote
	//FindQuotes(guildID string, search string) []Quote
}

func NewQuoteStore(uri string) (QuoteStore, error) {

	return newfileQuoteStore(uri)
}
