# go-metrics

Simplified metrics tracking system for golang to output in JSON format. This system is based on the library provided by [rcrowley](https://github.com/rcrowley/go-metrics)

Go Doc: <http://godoc.org/github.com/KyleLavorato/go-metrics>

## Supported Metric Types

* Counter - A basic integer value that can be set, incremented, or decremented.
* Json - A metric that will hold produced JSON values. Useful for aggregating previously output metrics into a single registry.
* Meter - A integer value that tracks past values applied. It will track the `count` of marked values, the `last` value marked and the `mean` or average value marked.
* Registry - The container that holds all metrics
* Slice - A generic array that will hold multiple metric entries of the same type
* Text - A simple string value that can be set or appended.
* Timer - An integer value that will track the number of seconds it executes for. It tracks the `count` of runs, a []integer set of `execution` times, the `lastValue`, as well as `max`, `mean`, and `min`.

## Basic Operations

### Metric Registration

The orchestrator of the metrics library is a `registry` object. This is a "map" that holds the values of all the metrics for the current session. There is a global registry that can always be used named `metrics.DefaultRegistry`. Users can also create their own master registry using the `metrics.NewRegistry()` function and managing the returned object.

Once a metric is created, it must be registered to a `registry` object. The register action will link its value to a metrics set and allow it to be tracked and eventually output. When registering the metric, a unique name will be used to map it into the registry. This will be the name the value is output with, and what it name can be used to retrieve it.

Additional registries can be created and registered to the master or other parent registries to create a nesting system. There is no limit to the amount of nesting that can be applied.

Each metric type supports two registration methods:

#### Register an Existing Metric
```go
c := metrics.NewCounter()
metrics.DefaultRegistry.Register("foo", c)
```
#### Create and Register Together
```go
c := metrics.NewRegisteredCounter("foo", registry)
```

### Metric Retrieval

The variable holding the reference to a metric does not need to be passed from function-to-function to be modified. The registry object supports a `Get()` function to retrieve the metric reference based on its registered name. 

When performing a `Get()` a type cast is required, as this call does not know the type of the metric that will be returned.

```go
metric := registry.Get("foo").(Counter)
```

### Output Metrics

To get the value of all the metrics contained a registry, simply call the `registry.GetAllJson()` function. This will dump the entire contents of the registry into JSON format to a byte variable.

The output can then be used as desired. It is recommended to forward this to an external metrics backend. Since JSON is a universal format, any number of different backend services can be used. 

The library currently supports output only in the JSON format.

## Installation

```sh
go get github.com/KyleLavorato/go-metrics
```

## Usage

Sample program that creates one of every metric type

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
)

type NestedData struct {
    ValueOne int
    ValueTwo int
}

type ExistingJson struct {
    Sample NestedData
    Data   NestedData
    Value  int
}

func Example() {
    registry := metrics.NewRegistry()

    // Example Counter
    c := metrics.NewRegisteredCounter("foo", registry)
    c.Inc(3)
    c.Dec(1)
    c.Inc(6)
    c2 := metrics.NewRegisteredCounter("bar", registry)
    c2.Inc(31)
    c2.Dec(14)
    c2.Inc(66)

    cx := registry.Get("foo").(Counter)
    cx.Inc(1)

    // Example Text
    t1 := metrics.NewRegisteredText("hello", registry)
    t1.Set("Error: ")
    t1.Append("did not hello world")

    // Example Meter
    m1 := metrics.NewRegisteredMeter("world", registry)
    m1.Mark(50)
    m1.Mark(100)

    // Example Timer
    q1 := metrics.NewRegisteredTimer("golang", registry)
    q1.Start()
    time.Sleep(time.Second * 1)
    q1.Stop()
    q1.Time(func() { time.Sleep(time.Second * 4) })

    // Example of Nesting Metrics
    nestedRegistry := metrics.NewRegistry()
    metrics.NewRegisteredText("msg", nestedRegistry).Set("This is a nested registry")
    metrics.NewRegisteredCounter("count", nestedRegistry).Inc(1996)
    registry.Register("registry2", nestedRegistry)

    // Example Slice
    s := metrics.NewSlice()
    for i := 0; i < 5; i++ {
        p := metrics.NewRegistry()
        metrics.NewRegisteredText("msg", p).Set("This is a slice entry")
        metrics.NewRegisteredCounter("count", p).Inc(int64(i))
        s.Append(p)
    }
    registry.Register("sliceReg", s)

    // Example JSON
    data := ExistingJson{
        Sample: NestedData{
            ValueOne: 1,
            ValueTwo: 2,
        },
        Data: NestedData{
            ValueOne: 3,
            ValueTwo: 4,
        },
        Value: 5,
    }
    raw, _ := json.Marshal(&data)
    metrics.NewRegisteredJson("ExistingJson", registry).Set(raw)

    // Example Output
    js, err := registry.GetAllJson()
    if err != nil {
        panic(err)
    }
    fmt.Println(string(js))
}
```

The expected output of the program is:

```json
{
  "ExistingJson": {
    "Sample": {
      "ValueOne": 1,
      "ValueTwo": 2
    },
    "Data": {
      "ValueOne": 3,
      "ValueTwo": 4
    },
    "Value": 5
  },
  "bar": 83,
  "foo": 9,
  "golang": {
    "count": 2,
    "executions": [
      1,
      4
    ],
    "lastValue": 4,
    "max": 4,
    "mean": 2.5,
    "min": 1
  },
  "hello": "Error: did not hello world",
  "registry2": {
    "count": 1996,
    "msg": "This is a nested registry"
  },
  "sliceReg": [
    {
      "count": 0,
      "msg": "This is a slice entry"
    },
    {
      "count": 1,
      "msg": "This is a slice entry"
    },
    {
      "count": 2,
      "msg": "This is a slice entry"
    },
    {
      "count": 3,
      "msg": "This is a slice entry"
    },
    {
      "count": 4,
      "msg": "This is a slice entry"
    }
  ],
  "world": {
    "count": 2,
    "lastValue": 100,
    "mean": 75
  }
}
```

## Additional Information

### Multi-Threading

The `Register()` function is not thread safe. If using a parallel processing model, ensure that registration calls are wrapped in a mutex for that registry.