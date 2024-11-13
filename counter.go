package metrics

import (
	"os"
	"sync/atomic"
)

// Counters hold an int64 value that can be incremented and decremented.
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Set(int64)
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	return &StandardCounter{0}
}

// NewRegisteredCounter constructs and registers a new StandardCounter.
func NewRegisteredCounter(name string, r Registry) Counter {
	c := NewCounter()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, c)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return c
}

// GetCounter returns an existing Counter
func GetCounter(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Counter)
}

// StandardCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardCounter struct {
	count int64
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}

// Set changes the counter to the given value.
func (c *StandardCounter) Set(i int64) {
	atomic.StoreInt64(&c.count, i)
}
