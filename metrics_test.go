// nolint
package metrics

import (
	"fmt"
	"time"
)

func Example() {
	registry := NewRegistry()

	c := NewRegisteredCounter("foo", registry)
	c.Inc(3)
	c.Dec(1)
	c.Inc(6)
	c2 := NewRegisteredCounter("bar", registry)
	c2.Inc(31)
	c2.Dec(14)
	c2.Inc(66)

	cx := registry.Get("foo").(Counter)
	cx.Inc(1)

	t1 := NewRegisteredText("hello", registry)
	t1.Set("Error: ")
	t1.Append("did not hello world")

	m1 := NewRegisteredMeter("world", registry)
	m1.Mark(50)
	m1.Mark(100)

	q1 := NewRegisteredTimer("golang", registry)
	q1.Start()
	time.Sleep(time.Second * 1)
	q1.Stop()
	q1.Time(func() { time.Sleep(time.Second * 4) })

	nestedRegistry := NewRegistry()
	NewRegisteredText("msg", nestedRegistry).Set("This is a nested registry")
	NewRegisteredCounter("count", nestedRegistry).Inc(1996)
	registry.Register("registry2", nestedRegistry)

	js, err := registry.GetAllJson()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(js))
	// Output: {"bar":83,"foo":9,"golang":{"count":2,"executions":[1,4],"lastValue":4,"max":4,"mean":2.5,"min":1},"hello":"Error: did not hello world","registry2":{"count":1996,"msg":"This is a nested registry"},"world":{"count":2,"lastValue":100,"mean":75}}
}
