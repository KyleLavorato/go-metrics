package metrics

import (
	"math"
	"math/rand"
	"testing"
)

func randomInt64() int64 {
	return int64(rand.Intn(100))
}

func BenchmarkMeter(b *testing.B) {
	m := NewMeter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Mark(randomInt64())
	}
}

func TestMeterCount(t *testing.T) {
	m := NewMeter()
	m.Mark(randomInt64())
	m.Mark(randomInt64())
	m.Mark(randomInt64())
	if count := m.Count(); 3 != count {
		t.Errorf("m.Count(): 3 != %v\n", count)
	}
}

func TestMeterMean(t *testing.T) {
	m := NewMeter()
	var total int64
	for i := 0; i < 10; i++ {
		val := randomInt64()
		m.Mark(val)
		total += val
	}
	mean := math.Float64frombits(math.Float64bits(float64(total) / float64(10)))
	if calc := m.RateMean(); calc != mean {
		t.Errorf("m.Mean(): %f != %f", calc, mean)
	}
}

func TestMeterDec2(t *testing.T) {
	m := NewMeter()
	var val int64
	for i := 0; i < 10; i++ {
		val = randomInt64()
		m.Mark(val)
	}
	if last := m.LastValue(); val != last {
		t.Errorf("m.LastValue(): val != %d", last)
	}
}

func TestMeterSnapshot(t *testing.T) {
	m := NewMeter()
	m.Mark(1)
	if snapshot := m.Snapshot(); m.RateMean() != snapshot.RateMean() {
		t.Fatal(snapshot)
	}
}

func TestMeterZero(t *testing.T) {
	m := NewMeter()
	if count := m.Count(); 0 != count {
		t.Errorf("m.Count(): 0 != %v\n", count)
	}
}

func TestGetMeter(t *testing.T) {
	r := NewRegistry()
	NewRegisteredMeter("foo", r).Mark(47)
	if m := GetMeter("foo", r); 47 != m.LastValue() || 47 != m.RateMean() {
		t.Fatal(m)
	}
}
