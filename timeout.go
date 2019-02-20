package tw

import (
	"container/list"
	"sync"
)

const (
	// states of a timeout
	timeoutInactive = iota
	timeoutExpired
	timeoutActive
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