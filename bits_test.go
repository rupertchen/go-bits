package bits

import "testing"

func TestNewBitmap(t *testing.T) {
	var tests = []struct{ capacity, storeSize int }{
		{0, 0},
		{1, 1},
		{64, 1},
		{65, 2},
		{128, 2},
		{129, 3},
	}

	for _, test := range (tests) {
		var b = NewBitmap(test.capacity)
		if b.Capacity() != test.capacity {
			t.Errorf("Expected %d capacity, got %d", test.capacity, b.Capacity())
		}
		if len(b.store) != test.storeSize {
			t.Errorf("Expected %d store size when capacit is %d, got %d", test.storeSize, test.capacity, len(b.store))
		}
	}
}

func TestBitmap_Get(t *testing.T) {

	var blockEquals = func(expected, actual Block) {
		if expected != actual {
			t.Errorf("Expected 0x%016X, got 0x%016X", expected, actual)
		}
	}

	var b = NewBitmap(64)
	blockEquals(0, b.Get(0, 0))
	blockEquals(0, b.Get(0, 1))

	b.store[0] = 0x8000000000000001
	blockEquals(0, b.Get(0, 0))
	blockEquals(0x1, b.Get(0, 1))
	blockEquals(0x1, b.Get(0, 2))
	blockEquals(0x1, b.Get(63, 1))
	blockEquals(0x2, b.Get(62, 2))

	var b2 = NewBitmap(128)
	b2.store[0] = 0xffffffffaaaaaaaa
	b2.store[1] = 0x5555555500000000
}
