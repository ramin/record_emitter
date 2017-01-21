package record_emitter

import (
	"io"
)

var (
	rByte byte = 13 // byte that corresponds to the '\r' rune.
	nByte byte = 10 // byte that corresponds to the '\n' rune.
)

type crlfReader struct {
	r io.Reader
}

// New creates a new io.Reader that wraps r to convert Classic Mac(CR)
// line endings to Linux (LF) line endings.
func NewCRLFReader(r io.Reader) io.Reader {
	return &crlfReader{r: r}
}

// Read replaces CR line endings in the source reader
// with LF line  endings.
func (r crlfReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	for i, b := range p {
		if b == rByte {
			p[i] = nByte
		}
	}
	return
}
