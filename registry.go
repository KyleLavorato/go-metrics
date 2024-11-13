package metrics

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Default registry is none is specified
var DefaultRegistry Registry = NewRegistry()

// DuplicateMetric is the error returned by Registry.Register when a metric
// already exists.  If you mean to Register that metric you must first
// Unregister the existing metric.
type DuplicateMetric string

func (err DuplicateMetric) Error() string {
	return fmt.Sprintf("duplicate metric: %s\n", string(err))
}

// The standard implementation of a Registry is a mutex-protected map
// of names to metrics.
type StandardRegistry struct {
	metrics map[string]interface{}
	mutex   sync.RWMutex
}

// A Registry holds references to a set of metrics by name and can iterate
// over them, calling callback functions provided by the user.
//
// This is an interface so as to encourage other structs to implement
// the Registry API as appropriate.
type Registry interface {
	// Call the given function for each registered metric.
	Each(func(string, interface{}))

	// Get the metric by the given name or nil if none is registered.
	Get(string) interface{}

	// Output the value of all metrics in JSON
	GetAllJson() ([]byte, error)

	// Register the given metric under the given name.
	Register(string, interface{}) error

	// Unregister the metric with the given name.
	Unregister(string)

	// Get the number of tracked metrics
	MetricCount() int
}

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface{})) {
	metrics := r.registered()
	for i := range metrics {
		kv := &metrics[i]
		f(kv.name, kv.value)
	}
}

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.metrics[name]
}

// Output the value of all registered metrics
func (r *StandardRegistry) serializeRegistry() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	data := make(map[string]interface{})
	r.Each(func(name string, i interface{}) {
		values := make(map[string]interface{})
		switch metric := i.(type) {
		case Counter:
			data[name] = metric.Count()
		case Meter:
			m := metric.Snapshot()
			values["count"] = m.Count()
			values["mean"] = m.RateMean()
			values["lastValue"] = m.LastValue()
			data[name] = values
		case Timer:
			t := metric
			values["count"] = t.Count()
			values["min"] = t.Min()
			values["max"] = t.Max()
			values["mean"] = t.Mean()
			values["lastValue"] = t.LastValue()
			values["executions"] = t.AllExecutions()
			data[name] = values
		case Histogram:
			h := metric
			ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
			values["count"] = h.Count()
			values["min"] = h.Min()
			values["max"] = h.Max()
			values["mean"] = h.Mean()
			values["stddev"] = h.StdDev()
			values["median"] = ps[0]
			values["75%"] = ps[1]
			values["95%"] = ps[2]
			values["99%"] = ps[3]
			values["99.9%"] = ps[4]
			data[name] = values
		case Text:
			data[name] = metric.Text()
		case Slice:
			slices := []interface{}{}
			for _, r := range metric.GetAll() {
				nestedReg := r.(*StandardRegistry)
				slices = append(slices, nestedReg.serializeRegistry())
			}
			data[name] = slices
		case Registry:
			nestedReg := metric.(*StandardRegistry)
			data[name] = nestedReg.serializeRegistry()
		}
	})
	return data
}

// Output the value of all registered metrics in JSON format
func (r *StandardRegistry) GetAllJson() ([]byte, error) {
	data := r.serializeRegistry()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (r *StandardRegistry) Register(name string, i interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.register(name, i)
}

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.metrics, name)
}

// Get the number of tracked metrics
func (r *StandardRegistry) MetricCount() int {
	return len(r.registered())
}

// Create a new registry.
func NewRegistry() Registry {
	return &StandardRegistry{metrics: make(map[string]interface{})}
}

func (r *StandardRegistry) register(name string, i interface{}) error {
	if _, ok := r.metrics[name]; ok {
		return DuplicateMetric(name)
	}
	switch i.(type) {
	case Counter, Text, Meter, Timer, Histogram, Registry, Slice:
		r.metrics[name] = i
	}
	return nil
}

type metricKV struct {
	name  string
	value interface{}
}

func (r *StandardRegistry) registered() []metricKV {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	metrics := make([]metricKV, 0, len(r.metrics))
	for name, i := range r.metrics {
		metrics = append(metrics, metricKV{
			name:  name,
			value: i,
		})
	}
	return metrics
}
