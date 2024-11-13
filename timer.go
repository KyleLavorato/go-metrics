package metrics

import (
	"math"
	"os"
	"sync"
	"time"
)

// Timers measures the time to complete a task and tracks the
// history of past executions
type Timer interface {
	Count() int64             // Number of timer executions
	Max() int64               // Highest recorded timer execution in seconds
	Mean() float64            // Mean recorded timer execution in seconds
	Min() int64               // Lowest recorded timer execution in seconds
	LastValue() float64       // The value from the most recent execution in seconds
	AllExecutions() []float64 // All past executions in seconds
	Start()                   // Record current time
	Stop()                    // Record duration since Start() call
	Time(func())              // Record duration to execute a function
}

// NewRegisteredTimer constructs and registers a new StandardTimer.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func NewRegisteredTimer(name string, r Registry) Timer {
	c := NewTimer()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, c)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return c
}

// NewTimer constructs a new StandardTimer using an exponentially-decaying
// sample with the same reservoir size and alpha as UNIX load averages.
// Be sure to call Cleanup() once the timer is of no use to allow for garbage collection.
func NewTimer() Timer {
	return &StandardTimer{
		histogram: NewHistogram(NewExpDecaySample(1028, 0.015)),
		meter:     NewMeter(),
	}
}

// GetTimer returns an existing Counter
func GetTimer(name string, r Registry) Timer {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Timer)
}

// StandardTimer is the standard implementation of a Timer and uses a Histogram
// and Meter.
type StandardTimer struct {
	histogram  Histogram
	meter      Meter
	executions []float64
	startTime  time.Time
	lastValue  float64
	mutex      sync.Mutex
}

// Count returns the number of events recorded.
func (t *StandardTimer) Count() int64 {
	return t.histogram.Count()
}

// Max returns the maximum value in the sample.
func (t *StandardTimer) Max() int64 {
	return t.histogram.Max()
}

// Mean returns the mean of the values in the sample.
func (t *StandardTimer) Mean() float64 {
	return t.histogram.Mean()
}

// Min returns the minimum value in the sample.
func (t *StandardTimer) Min() int64 {
	return t.histogram.Min()
}

// LastValue returns the value from the most recent execution.
func (t *StandardTimer) LastValue() float64 {
	return t.lastValue
}

// AllExecutions returns all the past executions
func (t *StandardTimer) AllExecutions() []float64 {
	return t.executions
}

// Record the current time to prepare for a Stop() call
func (t *StandardTimer) Start() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.startTime = time.Now()
}

// Record the duration of an event that started with a call to Start() and ends now.
func (t *StandardTimer) Stop() {
	t.update(time.Since(t.startTime))
}

// Record the duration of the execution of the given function.
func (t *StandardTimer) Time(f func()) {
	ts := time.Now()
	f()
	t.update(time.Since(ts))
}

// Record the duration of an event.
func (t *StandardTimer) update(d time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.histogram.Update(int64(d.Seconds()))
	t.meter.Mark(1)
	t.executions = append(t.executions, math.Round(d.Seconds()*100)/100)
	t.lastValue = math.Round(d.Seconds()*100) / 100
}
