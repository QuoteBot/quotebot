package pagination

import (
	"time"

	"github.com/QuoteBot/quotebot/pkg/datastorage"
)

type State struct {
	quotes      []datastorage.Quote
	curPage     int
	maxPage     int
	lastPageLen int
	lastSeen    time.Time
}

func NewState(quotes []datastorage.Quote) *State {

	l := len(quotes)
	maxPage := l / pageQuotes
	lastlen := l % pageQuotes
	if lastlen > 0 {
		maxPage++
	} else {
		lastlen = pageQuotes
	}

	return &State{
		quotes:      quotes,
		curPage:     1,
		maxPage:     maxPage,
		lastPageLen: lastlen,
		lastSeen:    time.Now(),
	}
}

//compute current page
func (state *State) GetCurrentPage() *Page {
	islast := state.curPage == state.maxPage
	isfirst := state.curPage == 1
	startPosition := (state.curPage - 1) * pageQuotes
	endPosition := startPosition + pageQuotes
	if islast {
		endPosition = startPosition + state.lastPageLen
	}

	return &Page{
		Values:  state.quotes[startPosition:endPosition],
		HasNext: !islast,
		HasPrev: !isfirst,
	}
}
