package tw

import (
	"log"
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
	buckets      []timeoutList
	freeBucket   []timeoutList
	bMask        uint64
	tChan        chan timeoutList
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
		buckets:      make([]timeoutList, defaultBucketSize),
		freeBucket:   make([]timeoutList, defaultBucketSize),
		bMask:        defaultBucketSize - 1,
	}
}

// Start starts the timeout wheel
func (c *Client) Start() {
	for c.state != stopped {
		switch c.state {
		case stopping:
			<-c.done
		case running:
			panic("Tried to start a running wheel")
		}
	}

	c.state = running
	c.done = make(chan struct{})
	c.tChan = make(chan timeoutList, defaultBucketSize)

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

func (c *Client) Schedule(d time.Duration, data []byte, cb func([]byte)) {
	if c.state != running {
		panic("system has stopped")
	}

	dTicks := (d + c.tickInterval - 1) / c.tickInterval
	deadline := atomic.LoadUint64(&c.ticks) + uint64(dTicks)

	t := &timeout{
		receiver: &OnTimeoutImpl{
			callback: cb,
		},
		userData: data,
		deadline: deadline,
	}
	c.Lock()
	defer c.Unlock()
	log.Printf("insert %v", deadline&c.bMask)
	b := c.buckets[deadline&c.bMask]
	// if the last tick has already passed the deadline, execute callback now
	if b.lastTick >= deadline {
		t.receiver.Callback(data)
	}
	// otherwise schedule timeout
	c.buckets[deadline&c.bMask].prepend(t)
}

func (c *Client) onTick() {
	var tl timeoutList
	ticker := time.NewTicker(c.tickInterval)
	for range ticker.C {
		atomic.AddUint64(&c.ticks, 1)
		log.Printf("tick: %v", c.ticks)
		c.Lock()
		if c.state != running {
			c.Unlock()
			break
		}

		bucket := c.buckets[c.ticks&c.bMask]
		bucket.lastTick = c.ticks
		t := bucket.head
		for t != nil {
			if t.deadline <= c.ticks {
				t.remove()
				tl.prepend(t)
			}
			t = t.next
		}
		c.Unlock()
		if tl.head == nil {
			continue
		}
		c.tChan <- tl
		c.Lock()
		tl.head = nil
		c.Unlock()
	}
	ticker.Stop()
}

// onExpire fires timeout callbacks
func (c *Client) onExpire() {
	for list := range c.tChan {
		t := list.head
		for t != nil {
			c.Lock()
			if t.receiver != nil {
				t.receiver.Callback(t.userData)
			}
			t = t.next
			c.Unlock()
		}
	}

	c.Lock()
	c.state = stopped
	c.Unlock()
	close(c.done)
}
