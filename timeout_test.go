package tw

import (
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

	first := tl.head
	counter := 0
	for first.next != nil {
		first = first.next
		counter++
	}
	if counter != 2 {
		t.Fatalf("incorrect linked list length! expected %d but got %d", 2, counter)
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
	counter := 0
	first := tl.head
	for first.next != nil {
		first = first.next
		counter++
	}
	if counter != 1 {
		t.Fatalf("incorrect linked list length! expected %d but got %d", 1, counter)
	}
}
