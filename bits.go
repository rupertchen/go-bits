// Package bits is a set of utilities for working with sequences of bits.
package bits

import "math"

// Bitmap represents a fixed-length, sequence of bits.
type Bitmap struct {
	cap   int
	store []uint64
}

const bitsPerBlock = 64

func NewBitmap(capacity int) *Bitmap {
	var arraySize = int(math.Ceil(float64(capacity) / bitsPerBlock))
	return &Bitmap{
		cap:   capacity,
		store: make([]uint64, arraySize),
	}
}

// Get returns a Block of bits.
func (b *Bitmap) Get(index, length uint) Block {
	if length == 0 {
		return 0
	}

	// TODO: Assume <= 64 bits
	var buf = b.store[0] >> index
	var mask uint64 = 0xFFFFFFFFFFFFFFFF >> (bitsPerBlock - length)
	return Block(buf & mask)
}

func (b *Bitmap) Capacity() int {
	return b.cap
}

// Block contains a sequence of 0–64 bits. If fewer than 64 bits are needed,
// padding is applied. It is up to the caller to know how many bits are used.
type Block uint64
