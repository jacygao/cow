package tw

import "time"

const (
	defaultTickInterval = time.Millisecond

	// Client states
	stopped = iota
	stopping
	running
)

type Config struct {
	tickInterval time.Duration
}

func DefaultConfig() Config {
	return Config{
		tickInterval: defaultTickInterval,
	}
}

type Option func(*Config)

type Client struct {
	ticks        uint64
	tickInterval time.Duration
	state        int
	done         chan struct{}
}

func New(options ...Option) *Client {
	conf := DefaultConfig()
	for _, f := range options {
		f(&conf)
	}
	return &Client{
		ticks:        0,
		tickInterval: conf.tickInterval,
		state:        stopped,
	}
}

func (c *Client) Start() {
	c.state = running
}

func (c *Client) Stop() {
	if c.state == running {
		c.state = stopping
	}
	<- c.done
}

func (c *Client) onTick() *time.Ticker {
	ticker := time.NewTicker(c.tickInterval)
	return ticker
}

func (c *Client) onExpire() {
	c.state = stopped
	close(c.done)
}