package tw

import (
	"sync"
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
	sync.Mutex
	ticks        uint64
	tickInterval time.Duration
	state        int
	buckets      []*timeouts
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
		buckets:      make([]*timeouts, defaultBucketSize),
		bMask:        defaultBucketSize - 1,
	}
}

func (c *Client) Start() {
	c.Lock()
	defer c.Unlock()

	for c.state != stopped {
		switch c.state {
		case stopping:
			c.Unlock()
			<-c.done
			c.Lock()
		case running:
			panic("Tried to start a running tw")
		}
	}

	c.state = running
	c.done = make(chan struct{})
	c.tChan = make(chan timeouts, defaultBucketSize)

	go c.onTick()
	go c.onExpire()
}

func (c *Client) Stop() {
	c.Lock()
	if c.state == running {
		c.state = stopping
		close(c.tChan)
	}
	c.Unlock()
	<-c.done
}

func (c *Client) Schedule(d time.Duration, cb func()) {
	dTicks := (d + c.tickInterval - 1) / c.tickInterval
	deadline := atomic.LoadUint64(&c.ticks) + uint64(dTicks)
	t := &timeout {
		callback: OnTimeoutImpl{},
		userData: cb,
		deadline: deadline,
	}
	c.Lock()
	defer c.Unlock()
	if c.buckets[deadline&c.bMask] == nil {
		c.buckets[deadline&c.bMask] = newTimeouts()
	}
	c.buckets[deadline&c.bMask].push(t)
}

func (c *Client) onTick() {
	var ts timeouts
	ticker := time.NewTicker(c.tickInterval)
	for range ticker.C {
		atomic.AddUint64(&c.ticks, 1)
		c.Lock()
		if c.state != running {
			c.Unlock()
			break
		}
		bucket := c.buckets[c.ticks&c.bMask]
		if bucket != nil {
			ts.set(bucket.list)
		}
		c.Unlock()
		if ts.list == nil {
			continue
		}

		select {
		case c.tChan <- ts:
			ts.unset()
		default:
		}
	}
	ticker.Stop()
}

func (c *Client) onExpire() {
	for ts := range c.tChan {
		if ts.list != nil {
			ts.Lock()
			for t := ts.list.Front(); t != nil; t = t.Next() {
				if timeout, ok := t.Value.(*timeout); ok {
					timeout.callback.Callback(timeout.userData)
					c.Lock()
					c.buckets[timeout.deadline&c.bMask].pop(t)
					c.Unlock()
				}
			}
			ts.Unlock()
		}
	}
	c.Lock()
	c.state = stopped
	c.Unlock()
	close(c.done)
}