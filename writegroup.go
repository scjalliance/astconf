package astconf

import "io"

type writegroup struct {
	io.Writer
	written int
	err     error
}

func (wg *writegroup) Write(p []byte) {
	if wg.err != nil {
		return
	}
	n, err := wg.Writer.Write(p)
	wg.written += n
	wg.err = err
}

func (wg *writegroup) Result() (n int, err error) {
	return wg.written, wg.err
}

func (wg *writegroup) Err() (err error) {
	return wg.err
}

func (wg *writegroup) Reset() {
	wg.written = 0
	wg.err = nil
}
