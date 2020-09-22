package tw

import (
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cli := New()
	cli.Start()
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		fmt.Println("test")
	})
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		fmt.Println("test2")
	})
	cli.Schedule(2*time.Second, []byte("123"), func([]byte) {
		fmt.Println("test3")
	})
	time.Sleep(3 * time.Second)
	cli.Stop()
	t.Fatal("")
}
