package peekingReader

import (
	"bufio"
	"errors"
	"io"
)

type bufReader struct {
	r  io.Reader
	br *bufio.Reader
}

// NewBufReader returns a PeekingReader.Reader. It uses an underlying
// *bufio.Reader to implement its capability.
func NewBufReader(r io.Reader) Reader {
	return &bufReader{r, bufio.NewReader(r)}
}

func (b *bufReader) Peek(n int) ([]byte, error) {
	return b.br.Peek(n)
}

func (b *bufReader) ReadByte() (byte, error) {
	return b.br.ReadByte()
}

func (b *bufReader) ReadRune() (rune, int, error) {
	return b.br.ReadRune()
}

func (b *bufReader) ReadBytes(size int) ([]byte, error) {
	s := make([]byte, size)
	actual := b.br.Buffered()
	if actual < size { // pull directly from reader since buffer is too small
		buf := make([]byte, actual) // get rest of buffer
		l, err := b.br.Read(buf)
		if err != nil {
			return nil, err
		}
		if l != actual {
			return nil, errors.New("couldn't get all of remaining buffer")
		}
		copy(s[:actual], buf)

		for actual != size { // may need to read more than once to get full amount
			pull := make([]byte, size-actual) // pull what is left
			pSize, err := b.r.Read(pull)
			if err != nil {
				return nil, err
			}
			copy(s[actual:], pull)
			actual += pSize // bytes read from underlying reader + buffered bytes
		}
		b.br.Reset(b.r) // reset buffered reader state
		return s, nil
	}
	_, err := b.br.Read(s)
	return s, err // only return the actual valid number of bytes
}

func (b *bufReader) Read(result []byte) (int, error) {
	return b.br.Read(result)
}
