package cow

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cli := New()
	cli.Start()

	result := 0

	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})
	cli.Schedule(2*time.Second, []byte("123"), func([]byte) {
		result++
	})

	time.Sleep(3 * time.Second)
	cli.Stop()

	expectedCallbacks := 4
	if result != expectedCallbacks {
		t.Fatalf("expected %d callback to be executed but got %d", expectedCallbacks, result)
	}
}
