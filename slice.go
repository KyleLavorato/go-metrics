package metrics

import "os"

// Counters hold an int64 value that can be incremented and decremented.
type Slice interface {
	Append(Registry)
	GetAll() []Registry
	Clear()
}

// NewSlice constructs a new StandardSlice.
func NewSlice() Slice {
	return &StandardSlice{}
}

// NewRegisteredSlice constructs and registers a new StandardCounter.
func NewRegisteredSlice(name string, r Registry) Slice {
	s := NewSlice()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, s)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return s
}

// GetSlice returns an existing Counter
func GetSlice(name string, r Registry) Slice {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Slice)
}

// StandardSlice is the standard implementation of a Slice metrics set
type StandardSlice struct {
	data []Registry
}

// Append a registry onto the slice
func (s *StandardSlice) Append(r Registry) {
	s.data = append(s.data, r)
}

// Return all included registries
func (s *StandardSlice) GetAll() []Registry {
	return s.data
}

// Clear sets the slize to empty
func (s *StandardSlice) Clear() {
	s.data = []Registry{}
}
