package pagination

import (
	"context"
	"errors"
	"time"
)

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
	toAdd      *State
	id         string
}

type res struct {
	page *Page
	err  error
}

type mapPageHandler struct {
	states     map[string]*State
	context    context.Context
	actionChan chan action
	resChan    chan *res
}

func newMapPageHandler() PageHandler /*context.Context*/ {

	//ctx := context.WithCancel(context.Background(), nil)
	ph := &mapPageHandler{
		states:     make(map[string]*State),
		context:    context.Background(),
		actionChan: make(chan action),
		resChan:    make(chan *res),
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
func (pageHandler *mapPageHandler) _add(id string, state *State) {
	_, ok := pageHandler.states[id]
	if ok {
		pageHandler.resChan <- &res{
			page: nil,
			err:  errors.New("page already exists"),
		}
		return
	}
	pageHandler.states[id] = state
	pageHandler.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject nothing in res
func (pageHandler *mapPageHandler) _remove(id string) {
	delete(pageHandler.states, id)
}

//inject curent page as res after incrementing the page number if has next else return nil
func (pageHandler *mapPageHandler) _next(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	if state.curPage < state.maxPage {
		state.curPage++
	}
	pageHandler.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject curent page as res after decrementing the page number if has previous else return nil
func (pageHandler *mapPageHandler) _prev(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	if state.curPage > 0 {
		state.curPage--
	}
	pageHandler.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject curent page
func (pageHandler *mapPageHandler) _cur(id string) {
	state, ok := pageHandler.states[id]
	if !ok {
		pageHandler.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	pageHandler.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
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

func (pageHandler *mapPageHandler) GetCurrentPage(embedID string) (*Page, error) {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getCur}
	r := <-pageHandler.resChan
	return r.page, r.err
}
func (pageHandler *mapPageHandler) GetNextPage(embedID string) (*Page, error) {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getNext}
	r := <-pageHandler.resChan
	return r.page, r.err
}
func (pageHandler *mapPageHandler) GetPreviousPage(embedID string) (*Page, error) {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: getPrev}
	r := <-pageHandler.resChan
	return r.page, r.err
}
func (pageHandler *mapPageHandler) Add(embedID string, state *State) (*Page, error) {
	pageHandler.actionChan <- action{id: embedID, toAdd: state, typeAction: add}
	r := <-pageHandler.resChan
	return r.page, r.err
}
func (pageHandler *mapPageHandler) Delete(embedID string) {
	pageHandler.actionChan <- action{id: embedID, toAdd: nil, typeAction: remove}
}
