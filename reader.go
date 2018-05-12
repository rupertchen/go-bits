package bits

// Reader provides a convenient way to read bits, up to 64 at a time.
// Successive calls to the Read* functions will read the next sequence of bits
// and advance the Reader's position. It is the responsibility of the caller to
// not read beyond the available range.
//
// This is *not* an implementation of io.Reader.
type Reader struct {
	position uint
	src      *Bitmap
}

const bitsPerBool = 1

// NewReader returns a Reader that reads from src.
func NewReader(src *Bitmap) *Reader {
	return &Reader{position: 0, src: src}
}

// ReadBits returns a Block containing the next n bits.
func (r *Reader) ReadBits(n uint) Block {
	var val = r.src.Get(r.position, n)
	r.position += n
	return val
}

// ReadBool returns the next bit as a bool.
func (r *Reader) ReadBool() bool {
	return 1 == r.ReadBits(bitsPerBool)
}

// ReadByte returns the next eight bits as a byte.
func (r *Reader) ReadByte() byte {
	return byte(r.ReadBits(bitsPerByte))
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
