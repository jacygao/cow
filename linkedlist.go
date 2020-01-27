package tw

import (
	"container/list"
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
	list *list.List
}

type timeout struct {
	callback OnTimeout
	userData interface{}
	deadline uint64
}

func (ts *timeouts) push(t *timeout) {
	ts.list.PushBack(t)
}

func (ts *timeouts) pop(t *list.Element) {
	ts.list.Remove(t)
}

func (ts *timeouts) set(l *list.List) {
	if l != nil {
		if l.Len() > 0 {
			ts.list = l
		}
	}
}

func (ts *timeouts) unset() {
	ts.list = list.New()
}
