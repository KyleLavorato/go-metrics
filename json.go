package metrics

import (
	"encoding/json"
	"os"
)

// Json is a basic string message that can be set and appended to
type Json interface {
	Clear()
	Json() json.RawMessage // Get the current msg value
	Set(json.RawMessage)   // Set a new msg value
}

// NewCounter constructs a new StandardJson.
func NewJson() Json {
	return &StandardJson{}
}

// NewRegisteredCounter constructs and registers a new StandardJson.
func NewRegisteredJson(name string, r Registry) Json {
	j := NewJson()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, j)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return j
}

// GetJson returns an existing Counter
func GetJson(name string, r Registry) Json {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Json)
}

// StandardJson is the standard implementation of a Json value
type StandardJson struct {
	raw json.RawMessage
}

// Clear removes the current msg value
func (t *StandardJson) Clear() {
	t.raw = nil
}

// Json returns the msg value
func (t *StandardJson) Json() json.RawMessage {
	return t.raw
}

// Set changes the msg value to the indicated string
func (t *StandardJson) Set(j json.RawMessage) {
	t.raw = j
}
