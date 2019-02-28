package tw

import (
	"container/list"
	"sync"
)

type OnTimeout interface {
	Callback(userData interface{})
}

type OnTimeoutImpl struct {
}

func (i OnTimeoutImpl) Callback(payload interface{}) {
	payload.(func())()
}

func newTimeouts() *timeouts {
	return &timeouts{
		list: list.New(),
	}
}

type timeouts struct {
	sync.Mutex
	list *list.List
}

type timeout struct {
	callback OnTimeout
	userData interface{}
	deadline uint64
}

func (ts *timeouts) push(t *timeout) {
	ts.Lock()
	ts.list.PushBack(t)
	ts.Unlock()
}

func (ts *timeouts) pop(t *list.Element) {
	ts.Lock()
	ts.list.Remove(t)
	ts.Unlock()
}

func (ts *timeouts) set(l *list.List) {
	ts.Lock()
	if l != nil {
		if l.Len() > 0 {
			ts.list = l
		}
	}
	ts.Unlock()
}

func (ts *timeouts) unset() {
	ts.Lock()
	ts.list = list.New()
	ts.Unlock()
}