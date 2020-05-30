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

type mapPageManager struct {
	states     map[string]*State
	context    context.Context
	actionChan chan action
	resChan    chan *res
}

func newMapPageManager() PageManager /*context.Context*/ {

	//ctx := context.WithCancel(context.Background(), nil)
	ph := &mapPageManager{
		states:     make(map[string]*State),
		context:    context.Background(),
		actionChan: make(chan action),
		resChan:    make(chan *res),
	}
	go ph.behavior()
	return ph
}

//should be executed every gcTickDuration
func (PageManager *mapPageManager) _gc() {
	toDel := make([]string, 0, len(PageManager.states))
	now := time.Now()
	for id, s := range PageManager.states {
		if s.lastSeen.Add(timoutDuration).Before(now) {
			toDel = append(toDel, id)
		}
	}
	for _, id := range toDel {
		delete(PageManager.states, id)
	}
}

//inject current page as res
func (PageManager *mapPageManager) _add(id string, state *State) {
	_, ok := PageManager.states[id]
	if ok {
		PageManager.resChan <- &res{
			page: nil,
			err:  errors.New("page already exists"),
		}
		return
	}
	PageManager.states[id] = state
	PageManager.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject nothing in res
func (PageManager *mapPageManager) _remove(id string) {
	delete(PageManager.states, id)
}

//inject curent page as res after incrementing the page number if has next else return nil
func (PageManager *mapPageManager) _next(id string) {
	state, ok := PageManager.states[id]
	if !ok {
		PageManager.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	if state.curPage < state.maxPage {
		state.curPage++
	}
	PageManager.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject curent page as res after decrementing the page number if has previous else return nil
func (PageManager *mapPageManager) _prev(id string) {
	state, ok := PageManager.states[id]
	if !ok {
		PageManager.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	if state.curPage > 0 {
		state.curPage--
	}
	PageManager.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

//inject curent page
func (PageManager *mapPageManager) _cur(id string) {
	state, ok := PageManager.states[id]
	if !ok {
		PageManager.resChan <- &res{
			page: nil,
			err:  errors.New("page not found"),
		}
		return
	}
	PageManager.resChan <- &res{
		page: state.GetCurrentPage(),
		err:  nil,
	}
}

func (PageManager *mapPageManager) behavior() {
	gcTicker := time.NewTicker(gcTickDuration)

	for {
		select {
		case <-gcTicker.C:
			PageManager._gc()
		case action := <-PageManager.actionChan:
			switch action.typeAction {
			case add:
				PageManager._add(action.id, action.toAdd)
			case remove:
				PageManager._remove(action.id)
			case getCur:
				PageManager._cur(action.id)
			case getNext:
				PageManager._next(action.id)
			case getPrev:
				PageManager._prev(action.id)
			}

		case <-PageManager.context.Done():
			return
		}

	}
}

func (PageManager *mapPageManager) GetCurrentPage(embedID string) (*Page, error) {
	PageManager.actionChan <- action{id: embedID, toAdd: nil, typeAction: getCur}
	r := <-PageManager.resChan
	return r.page, r.err
}
func (PageManager *mapPageManager) GetNextPage(embedID string) (*Page, error) {
	PageManager.actionChan <- action{id: embedID, toAdd: nil, typeAction: getNext}
	r := <-PageManager.resChan
	return r.page, r.err
}
func (PageManager *mapPageManager) GetPreviousPage(embedID string) (*Page, error) {
	PageManager.actionChan <- action{id: embedID, toAdd: nil, typeAction: getPrev}
	r := <-PageManager.resChan
	return r.page, r.err
}
func (PageManager *mapPageManager) Add(embedID string, state *State) (*Page, error) {
	PageManager.actionChan <- action{id: embedID, toAdd: state, typeAction: add}
	r := <-PageManager.resChan
	return r.page, r.err
}
func (PageManager *mapPageManager) Delete(embedID string) {
	PageManager.actionChan <- action{id: embedID, toAdd: nil, typeAction: remove}
}
