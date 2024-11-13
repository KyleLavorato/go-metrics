package metrics

import (
	"math"
	"os"
	"sync/atomic"
	"time"
)

// Meters track set values and how the value changes over time
// each time it is recorded
type Meter interface {
	Count() int64      // Number of times the meter has saved a value
	Mark(int64)        // Set a value in the meter
	RateMean() float64 // Get the mean value from the meter
	LastValue() int64  // Get the most recently saved value
	Snapshot() Meter   // Save a snapshot of the current state
}

// NewMeter constructs a new StandardMeter and launches a goroutine.
// Be sure to call Stop() once the meter is of no use to allow for garbage collection.
func NewMeter() Meter {
	return &StandardMeter{
		snapshot:  &MeterSnapshot{},
		startTime: time.Now(),
	}
}

// NewMeter constructs and registers a new StandardMeter and launches a
// goroutine.
// Be sure to unregister the meter from the registry once it is of no use to
// allow for garbage collection.
func NewRegisteredMeter(name string, r Registry) Meter {
	c := NewMeter()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, c)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return c
}

// GetMeter returns an existing Counter
func GetMeter(name string, r Registry) Meter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Meter)
}

// MeterSnapshot is a read-only copy of another Meter.
type MeterSnapshot struct {
	count     int64
	value     int64
	rateMean  uint64
	lastValue int64
}

// Count returns the count of events at the time the snapshot was taken.
func (m *MeterSnapshot) Count() int64 { return m.count }

// Mark panics.
func (*MeterSnapshot) Mark(n int64) {
	panic("Mark called on a MeterSnapshot")
}

// RateMean returns the meter's mean rate of events per second at the time the
// snapshot was taken.
func (m *MeterSnapshot) RateMean() float64 { return math.Float64frombits(m.rateMean) }

// Snapshot returns the snapshot.
func (m *MeterSnapshot) Snapshot() Meter { return m }

// LastValue returns the last recorded value on the meter
func (m *MeterSnapshot) LastValue() int64 { return m.lastValue }

// StandardMeter is the standard implementation of a Meter.
type StandardMeter struct {
	snapshot  *MeterSnapshot
	startTime time.Time
}

// Count returns the number of events recorded.
func (m *StandardMeter) Count() int64 {
	return atomic.LoadInt64(&m.snapshot.count)
}

// Mark records the occurance of n events.
func (m *StandardMeter) Mark(n int64) {
	atomic.AddInt64(&m.snapshot.count, 1)
	atomic.AddInt64(&m.snapshot.value, n)
	atomic.StoreInt64(&m.snapshot.lastValue, n)

	m.updateSnapshot()
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeter) RateMean() float64 {
	return math.Float64frombits(atomic.LoadUint64(&m.snapshot.rateMean))
}

// LastValue returns the last recorded value on the meter
func (m *StandardMeter) LastValue() int64 { return atomic.LoadInt64(&m.snapshot.lastValue) }

// Snapshot returns a read-only copy of the meter.
func (m *StandardMeter) Snapshot() Meter {
	copiedSnapshot := MeterSnapshot{
		count:     atomic.LoadInt64(&m.snapshot.count),
		rateMean:  atomic.LoadUint64(&m.snapshot.rateMean),
		lastValue: atomic.LoadInt64(&m.snapshot.lastValue),
	}
	return &copiedSnapshot
}

func (m *StandardMeter) updateSnapshot() {
	rateMean := math.Float64bits(float64(atomic.LoadInt64(&m.snapshot.value)) / float64(m.Count()))

	atomic.StoreUint64(&m.snapshot.rateMean, rateMean)
}
