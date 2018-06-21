package astconf

import (
	"io"
	"reflect"
)

// Encoder encodes Go types as asterisk configuration data.
type Encoder struct {
	*encoderState // shared state among all encoders in a stack

	next *Encoder // the next encoder in the stack

	base   io.Writer // the underlying writer
	w      io.Writer // specialized writer wrapping base, may be the same as base
	indent int       // the current indent level of the encoder

	scratch encoderScratch
}

// NewEncoder returns a new Encoder that writes to the given io.Writer.
func NewEncoder(w io.Writer, options ...EncOpt) *Encoder {
	e := newEncoder(w)
	for _, option := range options {
		option(e)
	}
	return e
}

func newEncoder(w io.Writer) *Encoder {
	var e encoder

	// Initialize the stack
	for i := range e.stack {
		var next *Encoder
		if i+1 < len(e.stack) {
			next = &e.stack[i+1]
		}
		e.stack[i].init(&e.state, next)
	}

	// The root encoder is the first encoder in the stack
	root := &e.stack[0]

	// Prepare the root with a startup writer that will detect when the
	// encoder has been started
	root.base = w
	root.w = newStartupWriter(root)

	return root
}

// Printer returns a new Printer that writes to the same io.Writer as e.
func (e *Encoder) Printer() *Printer {
	return &Printer{
		Writer:    e.w,
		Indent:    e.indent,
		Alignment: e.alignment,
		Started:   e.started,
	}
}

// Encode encodes the given value as an asterisk configuration section.
func (e *Encoder) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	return valueEncoder(val)(val, e)
}

// init initializes a new Encoder.
func (e *Encoder) init(state *encoderState, next *Encoder) {
	e.encoderState = state
	e.next = next
}

// reset prepares an encoder for re-use
func (e *Encoder) reset() {
	e.indent = 0
}

// derive returns a child encoder that writes to e.w.
func (e *Encoder) derive() *Encoder {
	next := e.next
	if next == nil {
		next = &Encoder{}
		next.init(e.encoderState, nil)
	}
	next.base = e.w
	next.w = e.w
	next.indent = e.indent
	return next
}

// encoder is an internal struct that holds shared encoder state and
// a stack of encoders that can be used by marshaling functions.
type encoder struct {
	state encoderState
	stack [8]Encoder // Allocate an array of encoders as a contiguous memory block
}

// encoderState holds encoder state that is common to an encoder and
// all of its children.
type encoderState struct {
	alignment Alignment
	started   bool
}

// encoderScratch holds ephemeral encoder data used by encoding functions.
type encoderScratch struct {
	name    []byte
	sep     []byte // = or =>
	started bool
}
