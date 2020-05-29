package pagination

import (
	"context"
	"time"

	"github.com/QuoteBot/quotebot/pkg/datastorage"
)

type timedState struct {
	quotes      []datastorage.Quote
	curPage     int
	maxPage     int
	lastPageLen int
	lastSeen    time.Time
}

func fromQuotes(quotes []datastorage.Quote) *timedState {

	l := len(quotes)
	maxPage := l / pageQuotes
	lastlen := l % pageQuotes
	if lastlen > 0 {
		maxPage++
	} else {
		lastlen = pageQuotes
	}

	return &timedState{
		quotes:      quotes,
		curPage:     1,
		maxPage:     maxPage,
		lastPageLen: lastlen,
		lastSeen:    time.Now(),
	}
}

const gcTickDuration = 20 * time.Second
const timoutDuration = 2 * time.Minute

type actionType int

const (
	add     = actionType(0)
	remove  = actionType(1)
	getPrev = actionType(2)
	getNext = actionType(3)
	getCur  = actionType(4)
)

type action struct {
	typeAction actionType
	toAdd      *timedState
	id         string
}

type mapPageHandler struct {
	states     map[string]*timedState
	context    context.Context
	actionChan chan action
	resChan    chan *Page
}

func newMapPageHandler() PageHandler /*context.Context*/ {

	//ctx := context.WithCancel(context.Background(), nil)
	ph := &mapPageHandler{
		states:     make(map[string]*timedState),
		context:    context.Background(),
		actionChan: make(chan action),
		resChan:    make(chan *Page),
	}
	go ph.behavior()
	return ph
}

//should be executed every gcTickDuration
func (pageHandler *mapPageHandler) _gc() {
	toDel := make([]string, 0, len(pageHandler.states))
	now := time.Now()
	for id, s := range pageHandler.states {
		if s.lastSeen.Add(timoutDuration).Before(now) {
			toDel = append(toDel, id)
		}
	}
	for _, id := range toDel {
		delete(pageHandler.states, id)
	}
}

//inject current page as res
func (pageHandler *mapPageHandler) _add(id string, state *timedState) {
	_, ok := pageHandler.states[id]
	if ok {
		pageHandler.resChan <- nil
		return
	}
	pageHandler.states[id] = state
	res := cur(state)
	pageHandler.resChan <- res
}

//inject nothing in res
func (pageHandler *mapPageHandler) _remove(id string) {
	delete(pageHandler.states, id)
}

//inject curent page as res after incrementing the page number if has next else return nil
func (pageHandler *mapPageHandler) _next(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- nil
		return
	}
	if state.curPage > state.maxPage {
		state.curPage--
	}
	res := cur(state)
	pageHandler.resChan <- res
}

//inject curent page as res after decrementing the page number if has previous else return nil
func (pageHandler *mapPageHandler) _prev(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- nil
		return
	}
	if state.curPage > 0 {
		state.curPage--
	}
	res := cur(state)
	pageHandler.resChan <- res

}

//inject curent page
func (pageHandler *mapPageHandler) _cur(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- nil
		return
	}
	res := cur(state)
	pageHandler.resChan <- res
}

//compute current page
func cur(state *timedState) *Page {
	islast := state.curPage == state.maxPage
	isfirst := state.curPage == 1
	startPosition := (state.curPage - 1) * pageQuotes
	endPosition := startPosition + pageQuotes
	if islast {
		endPosition = startPosition + state.lastPageLen
	}

	return &Page{
		values:  state.quotes[startPosition:endPosition],
		hasNext: !islast,
		hasPrev: !isfirst,
	}
}

func (pageHandler *mapPageHandler) behavior() {
	gcTicker := time.NewTicker(gcTickDuration)

	for {
		select {
		case <-gcTicker.C:
			pageHandler._gc()
		case action := <-pageHandler.actionChan:
			switch action.typeAction {
			case add:
				pageHandler._add(action.id, action.toAdd)
			case remove:
				pageHandler._remove(action.id)
			case getCur:
				pageHandler._cur(action.id)
			case getNext:
				pageHandler._next(action.id)
			case getPrev:
				pageHandler._prev(action.id)
			}

		case <-pageHandler.context.Done():
			return
		}

	}
}

func (pageHandler *mapPageHandler) GetCurrentPage(embedID string) *Page {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getCur}
	return <-pageHandler.resChan
}
func (pageHandler *mapPageHandler) GetNextPage(embedID string) *Page {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getNext}
	return <-pageHandler.resChan
}
func (pageHandler *mapPageHandler) GetPreviousPage(embedID string) *Page {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getPrev}
	return <-pageHandler.resChan
}
func (pageHandler *mapPageHandler) Add(embedID string, quotes []datastorage.Quote) *Page {
	state := fromQuotes(quotes)
	pageHandler.actionChan <- action{id: embedID, toAdd: state, typeAction: add}
	return <-pageHandler.resChan
}
func (pageHandler *mapPageHandler) Delete(embedID string) {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: remove}
}
