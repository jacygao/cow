package tw

type OnTimeout interface {
	Callback(userData []byte)
}

type OnTimeoutImpl struct {
	callback func([]byte)
}

func (ot *OnTimeoutImpl) Callback(userData []byte) {
	ot.callback(userData)
}

type timeout struct {
	mtx      *lock
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
		t.prev = nil
		t.next = nil
		tl.head = t
	} else {
		t.prev = nil
		t.next = tl.head
		tl.head.prev = t
		tl.head = t
	}
}
