package metrics

import "os"

// Text is a basic string message that can be set and appended to
type Text interface {
	Clear()
	Text() string  // Get the current msg value
	Set(string)    // Set a new msg value
	Append(string) // Append to the end of the msg value
}

// NewCounter constructs a new StandardText.
func NewText() Text {
	return &StandardText{""}
}

// NewRegisteredCounter constructs and registers a new StandardText.
func NewRegisteredText(name string, r Registry) Text {
	t := NewText()
	if nil == r {
		r = DefaultRegistry
	}
	err := r.Register(name, t)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	return t
}

// GetText returns an existing Counter
func GetText(name string, r Registry) Text {
	if nil == r {
		r = DefaultRegistry
	}
	return r.Get(name).(Text)
}

// StandardText is the standard implementation of a text value
type StandardText struct {
	msg string
}

// Clear removes the current msg value
func (t *StandardText) Clear() {
	t.msg = ""
}

// Text returns the msg value
func (t *StandardText) Text() string {
	return t.msg
}

// Set changes the msg value to the indicated string
func (t *StandardText) Set(str string) {
	t.msg = str
}

// Append adds the indicated string to the end of the msg value
func (t *StandardText) Append(str string) {
	t.msg += str
}
