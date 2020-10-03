package tw

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
	cli.Schedule(2*time.Second, []byte("123"), func([]byte) {
		result++
	})

	time.Sleep(3 * time.Second)
	cli.Stop()

	if result != 2 {
		t.Fatalf("expected 2 callback to be executed but got %d", result)
	}
}