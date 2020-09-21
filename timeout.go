package tw

import (
	"sync"
	"unsafe"
)

const (
	cacheline = 64
)

type OnTimeout interface {
	Callback(userData []byte)
}

type OnTimeoutImpl struct {
	callback func([]byte)
}

func (ot *OnTimeoutImpl) Callback(userData []byte) {
	ot.callback(userData)
}

// sync.Mutex padded to a cache line to avoid false sharing
type mutex struct {
	sync.Mutex
	_ [cacheline - unsafe.Sizeof(sync.Mutex{})]byte
}

type timeout struct {
	mtx      *mutex
	receiver OnTimeout
	userData []byte
	deadline uint64
	next     *timeout
}

func (t *timeout) remove() {
	t.userData = t.next.userData
	t.next = t.next.next
}

type timeoutList struct {
	lastTick uint64
	head     *timeout
}

func (tl *timeoutList) prepend(t *timeout) {
	t.next = tl.head
	tl.head = t
}
