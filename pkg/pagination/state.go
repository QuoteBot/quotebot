package pagination

import (
	"time"

	"github.com/QuoteBot/quotebot/pkg/datastorage"
	"github.com/bwmarrin/discordgo"
)

type State struct {
	quotes      []datastorage.Quote
	curPage     int
	maxPage     int
	lastPageLen int
	lastSeen    time.Time
	author      *discordgo.User
	mentioned   *discordgo.User
}

func NewState(quotes []datastorage.Quote, author *discordgo.User, mentioned *discordgo.User) *State {

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
		author:      author,
		mentioned:   mentioned,
	}
}

//compute current page
func (state *State) GetCurrentPage() *Page {
	islast := state.curPage == state.maxPage
	isfirst := state.curPage == 0
	startPosition := (state.curPage) * pageQuotes
	endPosition := startPosition + pageQuotes
	if islast {
		endPosition = startPosition + state.lastPageLen
	}

	state.lastSeen = time.Now()

	return &Page{
		Values:    state.quotes[startPosition:endPosition],
		HasNext:   !islast,
		HasPrev:   !isfirst,
		Author:    state.author,
		Mentioned: state.mentioned,
	}
}
