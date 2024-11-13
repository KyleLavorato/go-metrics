package metrics

import (
	"fmt"
	"sync"
	"testing"
)

func BenchmarkRegistry(b *testing.B) {
	r := NewRegistry()
	err := r.Register("foo", NewCounter())
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Each(func(string, interface{}) {})
	}
}

func BenchmarkHugeRegistry(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 10000; i++ {
		err := r.Register(fmt.Sprintf("foo%07d", i), NewCounter())
		if err != nil {
			b.Fatal(err)
		}
	}
	v := make([]string, 10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := v[:0]
		r.Each(func(k string, _ interface{}) {
			v = append(v, k)
		})
	}
}

func BenchmarkRegistryParallel(b *testing.B) {
	r := NewRegistry()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewRegisteredCounter("foo", r)
		}
	})
}

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	err := r.Register("foo", NewCounter())
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if "foo" != name {
			t.Fatal(name)
		}
		if _, ok := iface.(Counter); !ok {
			t.Fatal(iface)
		}
	})
	if 1 != i {
		t.Fatal(i)
	}
	r.Unregister("foo")
	i = 0
	r.Each(func(string, interface{}) { i++ })
	if 0 != i {
		t.Fatal(i)
	}
}

func TestRegistryDuplicate(t *testing.T) {
	r := NewRegistry()
	if err := r.Register("foo", NewCounter()); nil != err {
		t.Fatal(err)
	}
	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if _, ok := iface.(Counter); !ok {
			t.Fatal(iface)
		}
	})
	if 1 != i {
		t.Fatal(i)
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	err := r.Register("foo", NewCounter())
	if err != nil {
		t.Fatal(err)
	}
	if count := r.Get("foo").(Counter).Count(); 0 != count {
		t.Fatal(count)
	}
	r.Get("foo").(Counter).Inc(1)
	if count := r.Get("foo").(Counter).Count(); 1 != count {
		t.Fatal(count)
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()

	// First metric wins with GetOrRegister
	_ = NewRegisteredCounter("foo", r)
	m := NewRegisteredMeter("foo", r)
	if _, ok := m.(Counter); ok {
		t.Fatal(m)
	}

	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if name != "foo" {
			t.Fatal(name)
		}
		if _, ok := iface.(Counter); !ok {
			t.Fatal(iface)
		}
	})
	if i != 1 {
		t.Fatal(i)
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()
	err := r.Register("foo", NewCounter())
	if err != nil {
		t.Fatal(err)
	}
	err = r.Register("bar", NewMeter())
	if err != nil {
		t.Fatal(err)
	}
	err = r.Register("baz", NewTimer())
	if err != nil {
		t.Fatal(err)
	}
	if r.MetricCount() != 3 {
		t.Errorf("Metric count: %d != 3", r.MetricCount())
	}
	r.Unregister("bar")
	r.Unregister("baz")
	if r.MetricCount() != 1 {
		t.Errorf("Metric count: %d != 1", r.MetricCount())
	}
	if _, ok := r.Get("foo").(Counter); !ok {
		t.Fatal("'foo' was removed but should exist")
	}
	if _, ok := r.Get("bar").(Meter); ok {
		t.Fatal("'bar' exists but should have been removed")
	}

}

func TestConcurrentRegistryAccess(t *testing.T) {
	r := NewRegistry()

	NewRegisteredCounter("foo", r)

	signalChan := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(dowork chan struct{}) {
			defer wg.Done()
		}(signalChan)
	}

	close(signalChan) // Closing will cause all go routines to execute at the same time
	wg.Wait()         // Wait for all go routines to do their work

	// At the end of the test we should still only have a single "foo" Counter
	i := 0
	r.Each(func(name string, iface interface{}) {
		i++
		if "foo" != name {
			t.Fatal(name)
		}
		if _, ok := iface.(Counter); !ok {
			t.Fatal(iface)
		}
	})
	if 1 != i {
		t.Fatal(i)
	}
	r.Unregister("foo")
	i = 0
	r.Each(func(string, interface{}) { i++ })
	if 0 != i {
		t.Fatal(i)
	}
}

// exercise race detector
func TestRegisterAndRegisteredConcurrency(t *testing.T) {
	r := NewRegistry()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(r Registry, wg *sync.WaitGroup) {
		defer wg.Done()
		r.Each(func(name string, iface interface{}) {
		})
	}(r, wg)
	err := r.Register("foo", NewCounter())
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}
