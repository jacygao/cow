package tw

import "sync"

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

type timeouts struct {
	sync.Mutex
	list []*timeout
}

type timeout struct {
	callback OnTimeout
	userData interface{}
	deadline uint64
}

func (ts *timeouts) prepend(t *timeout) {
	ts.Lock()
	ts.list = append(ts.list, t)
	ts.Unlock()
}

func (ts *timeouts) pop() {
	ts.Lock()
	ts.list = ts.list[:len(ts.list)-1]
	ts.Unlock()
}