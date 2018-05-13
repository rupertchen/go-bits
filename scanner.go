package bits

type Scanner struct {
	position uint
	bitmap   *Bitmap
}

const bitsPerBool = 1
const bitsPerByte = 8

func NewScanner(b *Bitmap) *Scanner {
	return &Scanner{position: 0, bitmap: b}
}

func (s *Scanner) ReadBits(length uint) Block {
	var val = s.bitmap.Get(s.position, length)
	s.position += length
	return val
}

func (s *Scanner) ReadBool() bool {
	return 1 == s.ReadBits(bitsPerBool)
}

func (s *Scanner) ReadByte() byte {
	return byte(s.ReadBits(bitsPerByte))
}
