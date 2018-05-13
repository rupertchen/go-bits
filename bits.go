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

// Get returns a Block of bits.
func (b *Bitmap) Get(index, length uint) Block {
	if length == 0 {
		return 0
	}

	// TODO: assume request does not span internal blocks
	var storeIndex = index / bitsPerBlock
	var buf = b.store[storeIndex] >> (bitsPerBlock - (index % bitsPerBlock) - length)
	var mask Block = 0xFFFFFFFFFFFFFFFF >> (bitsPerBlock - length)
	//fmt.Printf("Index, Length: %d, %d\n", index, length)
	//fmt.Printf("  s:\t0x%016X\n", b.store[0])
	//fmt.Printf("  buf:\t0x%016X\n", buf)
	//fmt.Printf("  mask:\t0x%016X\n", mask)
	return Block(buf & mask)
}

func (b *Bitmap) Capacity() int {
	return b.cap
}
