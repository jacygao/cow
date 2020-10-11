package cow

import (
	"reflect"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cli := New()
	cli.Start()

	result := 0

	// schedule an immediate callback
	cli.Schedule(0, []byte("123"), func([]byte) {
		result++
	})

	// schedule 3 callbacks in 1 second
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})
	cli.Schedule(1*time.Second, []byte("123"), func([]byte) {
		result++
	})

	// schedule callback in 2 second
	cli.Schedule(2*time.Second, []byte("123"), func([]byte) {
		result++
	})

	// Wait for all callbacks to be triggered
	time.Sleep(3 * time.Second)
	cli.Stop()

	expectedCallbacks := 5
	if result != expectedCallbacks {
		t.Fatalf("expected %d callback to be executed but got %d", expectedCallbacks, result)
	}
}

func TestOption(t *testing.T) {
	customVal := 2 * time.Second
	cli := New(WithTickInterval(customVal))

	if !reflect.DeepEqual(cli.tickInterval, customVal) {
		t.Fatalf("incorrect config values. expected %+v but got %+v", customVal, cli.tickInterval)
	}
}
