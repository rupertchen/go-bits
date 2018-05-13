package bits

// Reader provides a convenient way to read bits, up to 64 at a time.
// Successive calls to the Read* functions will read the next sequence of bits
// and advance the Reader's position.
//
// This is *not* an implementation of io.Reader.
type Reader struct {
	position uint
	bitmap   *Bitmap
}

const bitsPerBool = 1

func NewReader(b *Bitmap) *Reader {
	return &Reader{position: 0, bitmap: b}
}

func (r *Reader) ReadBits(length uint) Block {
	var val = r.bitmap.Get(r.position, length)
	r.position += length
	return val
}

func (r *Reader) ReadBool() bool {
	return 1 == r.ReadBits(bitsPerBool)
}

func (r *Reader) ReadByte() byte {
	return byte(r.ReadBits(bitsPerByte))
}
