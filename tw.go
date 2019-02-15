package tw

import (
	"sync/atomic"
	"time"
)

const (
	defaultTickInterval = time.Millisecond
	defaultBucketSize   = 2048

	// states of a Client
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
	buckets      []timeouts
	bMask        uint64
	tChan        chan timeouts
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
		buckets:      make([]timeouts, defaultBucketSize),
		bMask:        defaultBucketSize - 1,
		tChan:        make(chan timeouts, defaultBucketSize),
	}
}

func (c *Client) Start() {
	c.state = running
	c.done = make(chan struct{})

	go c.onTick()
	go c.onExpire()
}

func (c *Client) Stop() {
	if c.state == running {
		c.state = stopping
		close(c.tChan)
	}
	<-c.done
}

func (c *Client) Schedule(d time.Duration, cb func()) {
	dTicks := (d + c.tickInterval - 1) / c.tickInterval
	deadline := atomic.LoadUint64(&c.ticks) + uint64(dTicks)

	bucket := &c.buckets[deadline&c.bMask]
	t := &timeout{
		callback: OnTimeoutImpl{},
		userData: cb,
		deadline: deadline,
	}
	bucket.prepend(t)
}

func (c *Client) onTick() {
	var ts timeouts
	ticker := time.NewTicker(c.tickInterval)
	for range ticker.C {
		atomic.AddUint64(&c.ticks, 1)
		if c.state != running {
			break
		}
		bucket := &c.buckets[c.ticks&c.bMask]
		if len(bucket.list) > 0 {
			ts = *bucket
		}
		c.tChan <- ts
	}
	ticker.Stop()
}

func (c *Client) onExpire() {
	for ts := range c.tChan {
		for _, t := range ts.list {
			if t.callback != nil {
				t.callback.Callback(t.userData)
			}
		}
	}
	c.state = stopped
	close(c.done)
}

func (c *Client) lockAllBuckets() {
}

func (c *Client) unLockAllBuckets() {
}