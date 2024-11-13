package metrics

import (
	"math"
	"testing"
	"time"
)

func BenchmarkTimer(b *testing.B) {
	tm := NewTimer()
	b.ResetTimer()
	tm.Start()
	for i := 0; i < b.N; i++ {
		tm.Stop()
	}
}

func TestTimerCount(t *testing.T) {
	tm := NewTimer()
	tm.Time(func() { time.Sleep(time.Microsecond * 200) })
	tm.Time(func() { time.Sleep(time.Microsecond * 200) })
	tm.Time(func() { time.Sleep(time.Microsecond * 200) })
	if count := tm.Count(); 3 != count {
		t.Errorf("tm.Count(): 3 != %v\n", count)
	}
}

func TestTimerStartStop(t *testing.T) {
	tm := NewTimer()
	tm.Start()
	time.Sleep(time.Millisecond * 200)
	tm.Stop()
	if elapsed := tm.LastValue(); elapsed < float64(0.2) || elapsed > float64(0.203) {
		t.Errorf("tm.Mean(): %f != 0.2", elapsed)
	}
}

func TestTimerTime(t *testing.T) {
	tm := NewTimer()
	tm.Time(func() { time.Sleep(time.Millisecond * 200) })
	// Timer will never be exactly 0.2 due to timer overhead
	if elapsed := tm.LastValue(); elapsed < float64(0.2) || elapsed > float64(0.203) {
		t.Errorf("tm.LastValue(): %f != 0.2", elapsed)
	}
}

func TestTimerMetrics(t *testing.T) {
	tm := NewTimer()
	tm.Time(func() { time.Sleep(time.Second * 1) })
	tm.Time(func() { time.Sleep(time.Second * 2) })
	if mean := tm.Mean(); mean != float64(1.5) {
		t.Errorf("tm.Mean(): %f != 1.5", mean)
	}
	if max := tm.Max(); max != 2 {
		t.Errorf("tm.Max(): %d != 2", max)
	}
	if min := tm.Min(); min != 1 {
		t.Errorf("tm.Min(): %d != 1", min)
	}
}

func TestTimerZero(t *testing.T) {
	tm := NewTimer()
	if count := tm.Count(); 0 != count {
		t.Errorf("tm.Count(): 0 != %v\n", count)
	}
}

func TestTimerExecutions(t *testing.T) {
	tm := NewTimer()
	tm.Time(func() { time.Sleep(time.Microsecond * 200) })
	tm.Time(func() { time.Sleep(time.Microsecond * 200) })
	if ex := tm.AllExecutions(); 2 != len(ex) {
		t.Errorf("tm.Count(): 2 != %d", len(ex))
	}
}

func TestGetOrRegisterTimer(t *testing.T) {
	r := NewRegistry()
	NewRegisteredTimer("foo", r).Time(func() { time.Sleep(time.Millisecond * 200) })
	if tm := GetTimer("foo", r); 0.2 != math.Round(tm.LastValue()*10)/10 {
		t.Fatal(tm)
	}
}
