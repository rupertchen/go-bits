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
		if (b.Capacity() != test.capacity) {
			t.Errorf("Expected %d capacity, got %d", test.capacity, b.Capacity())
		}
		if len(b.store) != test.storeSize {
			t.Errorf("Expected %d store size, got %d", test.storeSize, len(b.store))
		}
	}
}
