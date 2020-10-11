package cow

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	defaultTickInterval = time.Millisecond
	defaultBucketSize   = 2048
	cacheline           = 64

	// states of a Client
	stopped = iota
	stopping
	running
)

// Config contains a list of values to initialise a Call Out Wheel.
type Config struct {
	tickInterval time.Duration
}

// DefaultConfig initialises configuration with default value.
func DefaultConfig() Config {
	return Config{
		tickInterval: defaultTickInterval,
	}
}

// Option overrides default configuration value with the input value.
type Option func(*Config)

// WithTickInterval sets the frequency of ticks.
func WithTickInterval(interval time.Duration) Option {
	return func(conf *Config) { conf.tickInterval = interval }
}

// Client defines a Wheel.
type Client struct {
	sync.Mutex
	ticks        uint64
	tickInterval time.Duration
	state        int
	locker       []lock
	buckets      []timeoutList
	bMask        uint64
	tChan        chan timeoutList
	done         chan struct{}
}

// sync.Mutex padded to a cache line to avoid false sharing
type lock struct {
	sync.Mutex
	_ [cacheline - unsafe.Sizeof(sync.Mutex{})]byte
}

// New returns a new Client instance.
func New(options ...Option) *Client {
	conf := DefaultConfig()
	for _, f := range options {
		f(&conf)
	}
	return &Client{
		ticks:        0,
		tickInterval: conf.tickInterval,
		state:        stopped,
		locker:       make([]lock, defaultBucketSize),
		buckets:      make([]timeoutList, defaultBucketSize),
		bMask:        defaultBucketSize - 1,
	}
}

// Start starts the call out wheel.
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

// Stop gracefully stops the wheel.
func (c *Client) Stop() {
	c.lockAll()

	if c.state == running {
		c.state = stopping
		close(c.tChan)
	}

	c.unlockAll()
	<-c.done
}

func (c *Client) leaseLock(deadline uint64) *lock {
	return &c.locker[deadline&c.bMask]
}

func (c *Client) lockAll() {
	for i := range c.locker {
		c.locker[i].Lock()
	}
}

func (c *Client) unlockAll() {
	for i := len(c.locker) - 1; i >= 0; i-- {
		c.locker[i].Unlock()
	}
}

// Schedule schedules a callback in a given time duration.
// The data is accesible by the callback function.
func (c *Client) Schedule(d time.Duration, data []byte, cb func([]byte)) bool {
	if c.state != running {
		panic("system has stopped")
	}

	dTicks := (d + c.tickInterval - 1) / c.tickInterval
	deadline := atomic.LoadUint64(&c.ticks) + uint64(dTicks)
	lock := c.leaseLock(deadline)
	lock.Lock()
	defer lock.Unlock()
	t := &timeout{}
	bucket := &c.buckets[deadline&c.bMask]
	t.deadline = deadline
	t.receiver = &OnTimeoutImpl{
		callback: cb,
	}
	t.userData = data
	// if the last tick has already passed the deadline, execute callback now
	if bucket.lastTick >= deadline {
		t.receiver.Callback(data)
		return true
	}
	// otherwise schedule timeout
	bucket.prepend(t)
	return true
}

// onTick runs on each tick
func (c *Client) onTick() {
	tl := timeoutList{}
	ticker := time.NewTicker(c.tickInterval)
	for range ticker.C {
		atomic.AddUint64(&c.ticks, 1)
		lock := c.leaseLock(c.ticks)
		lock.Lock()
		if c.state != running {
			lock.Unlock()
			break
		}
		bucket := &c.buckets[c.ticks&c.bMask]
		bucket.lastTick = c.ticks
		t := bucket.head
		for t != nil {
			next := t.next
			if t.deadline <= c.ticks {
				t.remove()
				tl.prepend(t)
			}
			t = next
		}
		lock.Unlock()
		if tl.head == nil {
			continue
		}
		c.tChan <- tl
		tl.head = nil
	}
	ticker.Stop()
}

// onExpire fires timeout callbacks.
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

	c.lockAll()
	c.state = stopped
	c.unlockAll()
	close(c.done)
}
