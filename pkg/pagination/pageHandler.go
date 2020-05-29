package pagination

import (
	"github.com/QuoteBot/quotebot/pkg/datastorage"
)

const pageQuotes = 5

//Page a page of selected quotes
type Page struct {
	Values  []datastorage.Quote
	HasNext bool
	HasPrev bool
}

//PageHandler interface for object handling pagination
type PageHandler interface {
	GetCurrentPage(embedID string) *Page
	GetNextPage(embedID string) *Page
	GetPreviousPage(embedID string) *Page
	Add(embedID string, state *State) *Page
	Delete(embedID string)
}

//NewPageHandler instantiate an object implementing the PageHandler interface
func NewPageHandler() PageHandler {
	return newMapPageHandler()
}
