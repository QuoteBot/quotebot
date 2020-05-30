package pagination

import (
	"github.com/QuoteBot/quotebot/pkg/datastorage"
	"github.com/bwmarrin/discordgo"
)

const pageQuotes = 5

//Page a page of selected quotes
type Page struct {
	Values    []datastorage.Quote
	HasNext   bool
	HasPrev   bool
	Author    *discordgo.User
	Mentioned *discordgo.User
}

//PageManager interface for object handling pagination
type PageManager interface {
	GetCurrentPage(embedID string) (*Page, error)
	GetNextPage(embedID string) (*Page, error)
	GetPreviousPage(embedID string) (*Page, error)
	Add(embedID string, state *State) (*Page, error)
	Delete(embedID string)
}

//NewPageManager instantiate an object implementing the PageManager interface
func NewPageManager() PageManager {
	return newMapPageManager()
}
