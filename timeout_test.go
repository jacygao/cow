package tw

import (
	"reflect"
	"testing"
)

func TestPrepend(t *testing.T) {
	t1 := &timeout{
		userData: []byte("test"),
	}
	t2 := &timeout{
		userData: []byte("test2"),
	}
	t3 := &timeout{
		userData: []byte("test3"),
	}

	tl := &timeoutList{}

	tl.prepend(t1)
	tl.prepend(t2)
	tl.prepend(t3)

	cout, ok := isLinkedListValid(tl)
	if !ok {
		t.Fatal("got invalid linked list")
	}
	if cout != 3 {
		t.Fatalf("incorrect linked list length! expected %d but got %d", 3, cout)
	}
}

func TestRemove(t *testing.T) {
	t1 := &timeout{
		userData: []byte("test"),
	}
	t2 := &timeout{
		userData: []byte("test2"),
	}
	t3 := &timeout{
		userData: []byte("test3"),
	}

	tl := &timeoutList{}

	tl.prepend(t1)
	tl.prepend(t2)
	tl.prepend(t3)

	t2.remove()
	cout, ok := isLinkedListValid(tl)
	if !ok {
		t.Fatal("got invalid linked list")
	}

	if cout != 2 {
		t.Fatalf("incorrect linked list length! expected %d but got %d", 1, cout)
	}
}

func isLinkedListValid(tl *timeoutList) (int, bool) {
	counter := 0
	ti := tl.head
	for ti != nil {
		if ti.next != nil {
			if !reflect.DeepEqual(ti, ti.next.prev) {
				return 0, false
			}
		}
		ti = ti.next
		counter++
	}
	return counter, true
}
