package bits

import (
	"fmt"

	"github.com/pkg/errors"
)

// Reader provides a convenient way to read bits, up to 64 at a time.
// Successive calls to the Read* functions will read the next sequence of bits
// and advance the Reader's position. It is the responsibility of the caller to
// not read beyond the available range.
//
// Once any Read* call returns an error, all subsequent calls on any Read*
// function will error, returning the same error as the first. This gives
// callers the option to "batch" read logic and periodically check for errors
// rather than checking it after every read.
//
// The internal position is not advanced when there is an error. Therefore, it
// is possible for NumUnread and HasUnread to return >0 or true, respectively,
// but be unable to read any more bits. This is an intentional design decision
// to prevent mis-aligned reads following a read failure.
//
// This is *not* an implementation of io.Reader.
type Reader struct {
	position uint
	src      *Bitmap
	Err      error
}

const bitsPerBool = 1

// NewReader returns a Reader that reads from src.
func NewReader(src *Bitmap) *Reader {
	return &Reader{position: 0, src: src}
}

// ReadBits returns a Block containing the next n bits.
func (r *Reader) ReadBits(n uint) (Block, error) {
	if r.Err != nil {
		return 0, r.Err
	}

	if b, err := r.src.Get(r.position, n); err != nil {
		r.Err = errors.WithMessage(err, fmt.Sprintf("read bits (index=%d, length=%d)", r.position, n))
		return 0, r.Err
	} else {
		r.position += n
		return b, nil
	}
}

// ReadBool returns the next bit as a bool.
func (r *Reader) ReadBool() (bool, error) {
	if r.Err != nil {
		return false, r.Err
	}

	if b, err := r.ReadBits(bitsPerBool); err != nil {
		r.Err = errors.WithMessage(err, "read bool")
		return false, r.Err
	} else {
		return 1 == b, nil
	}
}

// ReadByte returns the next eight bits as a byte.
func (r *Reader) ReadByte() (byte, error) {
	if r.Err != nil {
		return 0, r.Err
	}

	if b, err := r.ReadBits(bitsPerByte); err != nil {
		r.Err = errors.WithMessage(err, "read byte")
		return 0, r.Err
	} else {
		return byte(b), nil
	}
}

// Size returns the size of the underlying Bitmap.
func (r *Reader) Size() int {
	return r.src.Size()
}

// NumUnread returns the number of bits left to be read.
func (r *Reader) NumUnread() int {
	return int(uint(r.Size()) - r.position)
}

// HasUnread returns true iff there are unread bits.
func (r *Reader) HasUnread() bool {
	return uint(r.Size()) != r.position
}
