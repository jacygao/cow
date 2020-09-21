package tw

import (
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cli := New()
	cli.Start()
	cli.Schedule(2*time.Second, []byte("123"), func([]byte) {
		fmt.Println("test1")
	})
	time.Sleep(5 * time.Second)
	cli.Stop()
	t.Fatal("")
}
