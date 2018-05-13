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

	var b = NewBitmap(64)
	blockEquals(t, 0, b.Get(0, 0))
	blockEquals(t, 0, b.Get(0, 1))

	b.store[0] = 0x8000000000000001
	blockEquals(t, 0, b.Get(0, 0))
	blockEquals(t, 0x1, b.Get(0, 1))
	blockEquals(t, 0x2, b.Get(0, 2))
	blockEquals(t, 0x1, b.Get(63, 1))
	blockEquals(t, 0x1, b.Get(62, 2))

	var b2 = NewBitmap(128)
	b2.store[0] = 0xffffffffaaaaaaaa
	b2.store[1] = 0x55555555ffffffff

	// Access within first internal block
	blockEquals(t, 0xff, b2.Get(4, 8))
	blockEquals(t, 0xffaa, b2.Get(24, 16))
	blockEquals(t, 0x2a, b2.Get(58, 6))

	// Access within second internal block
	blockEquals(t, 0x55, b2.Get(64, 8))
	blockEquals(t, 0x55ff, b2.Get(88, 16))
	blockEquals(t, 0x3f, b2.Get(122, 6))

	// Access spanning internal blocks
	blockEquals(t, 0xa5, b2.Get(60, 8))
}

func blockEquals(t *testing.T, expected, actual Block) {
	if expected != actual {
		t.Errorf("Expected 0x%016X, got 0x%016X", expected, actual)
		t.FailNow()
	}
}
