package tw

import (
	"fmt"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cli := New()
	cli.Start()
	cli.Schedule(time.Second, func(){
		fmt.Println("test1")
	})
	cli.Schedule(time.Second, func(){
		fmt.Println("test2")
	})
	cli.Schedule(time.Second, func(){
		fmt.Println("test3")
	})
	cli.Schedule(2 * time.Second, func(){
		fmt.Println("test4")
	})
	cli.Schedule(4 * time.Second, func(){
		fmt.Println("test5")
	})
	time.Sleep(15 * time.Second)
	cli.Stop()
}