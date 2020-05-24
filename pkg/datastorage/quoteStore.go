package datastorage

import (
	"time"
)

type Quote struct {
	QuoteId   string    `json:"quoteID"`
	UserID    string    `json:"userID"`
	GuildID   string    `json:"guildId"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}

type QuoteStore interface {
	Save(quote *Quote) error
	Delete(quoteID string, userID string, guildID string) error
	//GetQuotesFromUser(userID string, guildID string) []Quote
	//GetAllQuotes(guildID string) []Quote
	//FindQuotesFromUser(userID string, guildID string, search string) []Quote
	//FindQuotes(guildID string, search string) []Quote
}

func NewQuoteStore(uri string) (QuoteStore, error) {

	return newfileQuoteStore(uri)
}
