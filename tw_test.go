package tw

import "testing"

func TestStartStop(t *testing.T) {
	cli := New()
	cli.Start()
	cli.Stop()
}