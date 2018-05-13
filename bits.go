// Package bits is a set of utilities for working with sequences of bits.
package bits

import "math"

// Block contains a sequence of 0–64 bits. If fewer than 64 bits are needed,
// padding is applied. It is up to the caller to know how many bits are used.
type Block uint64

// Bitmap represents a fixed-length, sequence of bits.
type Bitmap struct {
	cap   int
	store []Block
}

const bitsPerBlock = 64

func NewBitmap(capacity int) *Bitmap {
	var arraySize = int(math.Ceil(float64(capacity) / bitsPerBlock))
	return &Bitmap{
		cap:   capacity,
		store: make([]Block, arraySize),
	}
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

func (b *Bitmap) Capacity() int {
	// TODO: Rename to "Length"? "Capacity" implies this structure might be
	// writable.
	return b.cap
}
