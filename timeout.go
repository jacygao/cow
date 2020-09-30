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
	prev     *timeout
}

func (t *timeout) remove() {
	if t.prev != nil {
		t.prev.next = t.next
	}
	if t.next != nil {
		t.next.prev = t.prev
	}
}

type timeoutList struct {
	lastTick uint64
	head     *timeout
}

func (tl *timeoutList) prepend(t *timeout) {
	if tl.head == nil {
		tl.head = t
	} else {
		head := tl.head
		head.prev = t
		t.next = head
		t.prev = nil
		tl.head = t
	}
}
