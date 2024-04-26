package peekingReader

import "io"

var spaceChars = []byte{'\x00', '\t', '\f', ' ', '\n', '\r'}

// Reader is an io.Reader which can also Peek and Read n number of bytes.
// Also includes a convenience function to read one byte
type Reader interface {
	Peek(n int) ([]byte, error)
	ReadByte() (byte, error)
	ReadBytes(size int) ([]byte, error)
	ReadRune() (rune, int, error)
	io.Reader
}

// ReadUntil reads until a particular byte is reached
func ReadUntil(r Reader, endAt byte) ([]byte, error) {
	return ReadUntilAny(r, []byte{endAt})
}

// ReadUntilAny reads until any byte in the endAtAny array is reached
func ReadUntilAny(r Reader, endAtAny []byte) ([]byte, error) {
	var result []byte
	for {
		p, err := r.Peek(1)
		if err != nil {
			return result, err
		}
		for i := range endAtAny {
			if p[0] == endAtAny[i] {
				return result, err
			}
		}
		r.ReadByte() // move read pointer forward since we used this byte
		result = append(result, p[0])
	}
}

// SkipSpaces skips all subsequent space characters.
// This includes " ", "\t", "\f", "\n", "\r", "\x00"
func SkipSpaces(r Reader) error {
	_, err := SkipSubsequent(r, spaceChars)
	return err
}

// SkipSubsequent skips all consecutive characters in any order
func SkipSubsequent(r Reader, skip []byte) (bool, error) {
	var found bool
top:
	for {
		b, err := r.Peek(1) // check next byte
		if err != nil {
			return found, err
		}
		next := b[0]
		for i := range skip {
			if next == skip[i] { // found match, so do actual read to skip
				found = true
				r.ReadByte() // move read pointer forward since we used this byte
				continue top
			}
		}
		return found, nil
	}
}
