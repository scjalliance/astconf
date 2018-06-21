package astconf

import "io"

// startupWriter watches for the first successful write and updates
// the encoder's started state when it occurs.
//
// The startupWriter removes itself from the encoder once it has been
// applied.
type startupWriter struct {
	*Encoder
}

func newStartupWriter(e *Encoder) io.Writer {
	return startupWriter{e}
}

func (sw startupWriter) Write(p []byte) (n int, err error) {
	e := sw.Encoder
	if e.started {
		return e.w.Write(p)
	}
	n, err = e.base.Write(p)
	if n > 0 {
		e.started = true
		e.w = e.base
	}
	return n, err
}
