// Package bits is a set of utilities for working with sequences of bits.
package bits

// Block contains a sequence of 0–64 bits. If fewer than 64 bits are needed,
// padding is applied. It is up to the caller to know how many bits are used.
type Block uint64

// Bitmap represents a fixed-length, sequence of bits.
type Bitmap struct {
	size  int
	store []Block
}

const bitsPerBlock = bitsPerByte * bytesPerBlock
const bytesPerBlock = 8
const bitsPerByte = 8

// NewBitmap returns a new Bitmap of the specified size with all bits set to
// zero.
func NewBitmap(s int) *Bitmap {
	var arraySize = sizeRequired(s, bitsPerBlock)
	return &Bitmap{
		size:  s,
		store: make([]Block, arraySize),
	}
}

// NewBitmapFromBytes returns a new Bitmap from a slice of bytes.
func NewBitmapFromBytes(bytes []byte) *Bitmap {
	var numBytes = len(bytes)
	var storeSize = sizeRequired(numBytes, bytesPerBlock)
	var s = make([]Block, 0, storeSize)

	var buf Block
	var lastM = bytesPerBlock - 1
	for i, b := range bytes {
		var m = i % bytesPerBlock
		buf |= Block(b) << uint(bitsPerBlock-(m+1)*bitsPerByte)
		if lastM == m {
			s = append(s, buf)
			buf = 0
		}
	}
	if buf != 0 {
		s = append(s, buf)
	}
	return &Bitmap{size: numBytes, store: s}
}

// NewBitmapFromBlocks returns a new Bitmap from a slice of Blocks. This
// a convenient way to explicitly create Bitmaps in code, such as when writing
// tests.
func NewBitmapFromBlocks(blocks []Block) *Bitmap {
	return &Bitmap{size: len(blocks) * bitsPerBlock, store: blocks}
}

// sizeRequired returns the number of groups required to fit n items into
// groups of at most size g. This is equivalent to ceil(n / g) but avoids using
// package math, which depends on package unsafe.
func sizeRequired(n, g int) int {
	var s = n / g
	if 0 == (n % g) {
		return s
	}
	return s + 1
}

// Get returns a Block of bits representing the requested range of bits. The
// bits are right-aligned.
func (b *Bitmap) Get(index, length uint) Block {
	// TODO: Validate index and length

	if length == 0 {
		return 0
	}

	// This is the index of the internal block where the range begins.
	var storeIndex = index / bitsPerBlock

	// This is the index of the last bit relative to the internal block where
	// the range begins.
	var endBitIndex = index%bitsPerBlock + length - 1

	// Bit shift and merge the internal blocks containing the range.
	var buf Block
	if endBitIndex < bitsPerBlock {
		buf = b.store[storeIndex] >> (bitsPerBlock - endBitIndex - 1)
	} else {
		var overflow = endBitIndex - bitsPerBlock + 1
		buf = (b.store[storeIndex] << overflow) | (b.store[storeIndex+1] >> (bitsPerBlock - overflow))
	}

	var mask Block = 0xFFFFFFFFFFFFFFFF >> (bitsPerBlock - length)
	return Block(buf & mask)
}

func (b *Bitmap) Size() int {
	return b.size
}
