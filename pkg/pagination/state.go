package pagination

import (
	"sort"
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

//NewState build a new state
//
func NewState(quotes []datastorage.Quote, author *discordgo.User, mentioned *discordgo.User) *State {

	//sort quotes by score
	sortedQuotes := quotes[:]
	sort.Slice(sortedQuotes, func(i, j int) bool {
		if sortedQuotes[i].Score == sortedQuotes[j].Score {
			return sortedQuotes[i].Timestamp.Before(sortedQuotes[j].Timestamp)
		}
		return sortedQuotes[i].Score > sortedQuotes[j].Score
	})
	//define the number of pages and the size of the last page
	l := len(quotes)
	maxPage := l / pageQuotes
	lastlen := l % pageQuotes
	if lastlen > 0 {
		maxPage++
	} else {
		lastlen = pageQuotes
	}

	return &State{
		quotes:      sortedQuotes,
		curPage:     0,
		maxPage:     maxPage,
		lastPageLen: lastlen,
		lastSeen:    time.Now(),
		author:      author,
		mentioned:   mentioned,
	}
}

//GetCurrentPage compute current page
func (state *State) GetCurrentPage() *Page {
	islast := state.curPage == state.maxPage-1
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
