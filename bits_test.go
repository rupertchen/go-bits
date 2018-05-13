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

	for _, test:= range(tests) {
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
	// TODO: Set the underlying store directly as we don't have helpers yet.
	var b = NewBitmap(1)
	blockEquals(t, 0x0000, b.Get(0, 0))
	blockEquals(t, 0x0000, b.Get(0, 1))

	b.store[0] = 0x0001
	blockEquals(t, 0x0000, b.Get(0, 0))
	blockEquals(t, 0x0001, b.Get(0, 1))
}

func blockEquals(t *testing.T, expected, actual Block) {
	if expected != actual {
		t.Errorf("Expected 0x%04X, got 0x%04X", expected, actual)
	}
}
