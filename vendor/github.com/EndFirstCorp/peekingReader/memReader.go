package peekingReader

import (
	"io"
	"unicode/utf8"
)

type memReader struct {
	buf []byte
	i   int
}

// NewMemReader returns a PeekingReader.Reader. It uses an underlying
// []byte buffer to provide its capability
func NewMemReader(b []byte) Reader {
	return &memReader{b, 0}
}

func (b *memReader) Peek(n int) ([]byte, error) {
	if b.i+n > len(b.buf) {
		return nil, io.EOF
	}
	return b.buf[b.i : b.i+n], nil
}

func (b *memReader) ReadRune() (rune, int, error) {
	var size int
	var r rune
	if b.i+1 > len(b.buf) {
		return 0, 0, io.EOF
	}
	r, size = rune(b.buf[b.i]), 1
	if r >= utf8.RuneSelf {
		r, size = utf8.DecodeRune(b.buf[b.i : b.i+2])
	}
	b.i += size
	return r, size, nil
}

func (b *memReader) ReadByte() (byte, error) {
	if b.i+1 > len(b.buf) {
		return 0, io.EOF
	}
	v := b.buf[b.i]
	b.i++
	return v, nil
}

func (b *memReader) ReadBytes(size int) ([]byte, error) {
	if b.i+size > len(b.buf) {
		return nil, io.EOF
	}
	v := b.buf[b.i : b.i+size]
	b.i += size
	return v, nil
}

func (b *memReader) Read(result []byte) (int, error) {
	if b.i == len(b.buf) { // nothing more to read
		return 0, io.EOF
	}
	end := len(result) + b.i
	if end > len(b.buf) {
		end = len(b.buf)
	}
	copy(result, b.buf[b.i:end])
	size := end - b.i
	b.i = end
	return size, nil
}
