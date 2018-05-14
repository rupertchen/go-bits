// Package bits is a set of utilities for working with sequences of bits.
package bits

import "errors"

// Block contains a sequence of 0–64 bits. If fewer than 64 bits are needed,
// the sequence is right-aligned and padded with zeros on the left. It is up
// to the caller to interpret how many bits are used.
type Block uint64

// Bitmap represents a fixed-length, sequence of bits.
type Bitmap struct {
	size  int
	store []Block
}

const bitsPerByte = 8
const bytesPerBlock = 8
const bitsPerBlock = bitsPerByte * bytesPerBlock

// NewBitmap returns a new Bitmap from a slice of bytes. Passing the output of
// a decoder is the typical way of creating a Bitmap.
func NewBitmap(bytes []byte) *Bitmap {
	var numBytes = len(bytes)
	var storeSize = sizeRequired(numBytes, bytesPerBlock)
	var s = make([]Block, 0, storeSize)

	var buf Block
	var hasPartialBuf bool
	var lastM = bytesPerBlock - 1
	for i, b := range bytes {
		var m = i % bytesPerBlock
		buf |= Block(b) << uint(bitsPerBlock-(m+1)*bitsPerByte)
		if lastM == m {
			s = append(s, buf)
			buf = 0
			hasPartialBuf = false
		} else {
			hasPartialBuf = true
		}
	}
	if hasPartialBuf {
		s = append(s, buf)
	}
	return &Bitmap{size: numBytes * bitsPerByte, store: s}
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
// returned bits are right-aligned.
func (b *Bitmap) Get(index, length uint) Block {
	if index >= uint(b.size) {
		panic(errors.New("index out of range"))
	}

	// I can't think of a reasonable use case for explicitly reading zero bits,
	// but returning an empty Block should be the result. Perhaps there is
	// behavior will simplify a caller's logic in some edge case.
	if length == 0 {
		return 0
	}

	if length > bitsPerBlock {
		panic(errors.New("length out of range, [0-64]"))
	}

	if index+length-1 >= uint(b.size) {
		panic(errors.New("length extends beyond range"))
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

// Size returns the number of bits in the Bitmap.
func (b *Bitmap) Size() int {
	return b.size
}
