package bits

import "testing"

func TestSizeRequired(t *testing.T) {
	var tests = []struct{ s, n, g int }{
		{0, 0, 2},
		{1, 1, 2},
		{1, 2, 2},
		{2, 3, 2},
		{57, 396, 7},
	}
	for _, test := range tests {
		var actual = sizeRequired(test.n, test.g)
		if test.s != actual {
			t.Errorf("Expected %d, got %d", test.s, actual)
		}
	}
}

func TestNewBitmap(t *testing.T) {
	// Load < 64 bits, all zeros
	var bmp0 = NewBitmap(make([]byte, 1))
	blockEquals(t, 0x0, bmp0.MustGet(0, 8))

	// Load < 64 bits
	var bmp1 = NewBitmap([]byte{0xec, 0xbc, 0x65, 0xde, 0x67})
	blockEquals(t, 0xecbc65de67, bmp1.MustGet(0, 40))

	// Load 64 bits
	var bmp2 = NewBitmap([]byte{0x62, 0x91, 0x07, 0xa1, 0x80, 0xbf, 0x7c, 0xe9})
	blockEquals(t, 0x629107a180bf7ce9, bmp2.MustGet(0, 64))

	// Load > 64 bits
	var bmp3 = NewBitmap([]byte{0x01, 0xba, 0x1c, 0xeb, 0xcd, 0xce, 0x33, 0x63, 0xd9, 0xb7, 0x8d, 0x5c})
	blockEquals(t, 0x01ba1cebcdce3363, bmp3.MustGet(0, 64))
	blockEquals(t, 0x0000000000d9b78d, bmp3.MustGet(64, 24))
}

func TestBitmap_MustGet(t *testing.T) {
	var bmp1 = NewBitmapFromBlocks([]Block{0x8000000000000001})
	blockEquals(t, 0, bmp1.MustGet(0, 0))
	blockEquals(t, 0x1, bmp1.MustGet(0, 1))
	blockEquals(t, 0x2, bmp1.MustGet(0, 2))
	blockEquals(t, 0x1, bmp1.MustGet(63, 1))
	blockEquals(t, 0x1, bmp1.MustGet(62, 2))

	var bmp2 = NewBitmapFromBlocks([]Block{
		0xffffffffaaaaaaaa,
		0x55555555ffffffff,
	})

	// Access within first internal block
	blockEquals(t, 0xff, bmp2.MustGet(4, 8))
	blockEquals(t, 0xffaa, bmp2.MustGet(24, 16))
	blockEquals(t, 0x2a, bmp2.MustGet(58, 6))

	// Access within second internal block
	blockEquals(t, 0x55, bmp2.MustGet(64, 8))
	blockEquals(t, 0x55ff, bmp2.MustGet(88, 16))
	blockEquals(t, 0x3f, bmp2.MustGet(122, 6))

	// Access spanning internal blocks
	blockEquals(t, 0xa5, bmp2.MustGet(60, 8))
}

func TestBitmap_MustGetPanicsOnOutOfRange(t *testing.T) {
	var bmp = NewBitmapFromBlocks([]Block{0, 0})
	assertPanic(t, func() { bmp.MustGet(128, 0) }, "bits: index out of range")
	assertPanic(t, func() { bmp.MustGet(128, 1) }, "bits: index out of range")
	assertPanic(t, func() { bmp.MustGet(0, 65) }, "bits: length out of range, [0-64]")
	assertPanic(t, func() { bmp.MustGet(120, 64) }, "bits: length extends beyond range")
}

func TestBitmap_GetPanicsOnOutOfRange(t *testing.T) {
	var bmp = NewBitmapFromBlocks([]Block{0, 0})

	var testCases = []struct {
		index, length uint
		expected      string
	}{
		{128, 0, "bits: index out of range"},
		{128, 1, "bits: index out of range"},
		{0, 65, "bits: length out of range, [0-64]"},
		{120, 64, "bits: length extends beyond range"},
	}

	for _, c := range testCases {
		var _, err = bmp.Get(c.index, c.length)
		if err == nil {
			t.Errorf("Expected error")
			t.FailNow()
		}
		if err.Error() != c.expected {
			t.Errorf(`Expect "%s", got "%s"`, c.expected, err.Error())
			t.FailNow()
		}
	}
}

func blockEquals(t *testing.T, expected, actual Block) {
	if expected != actual {
		t.Errorf("Expected 0x%016x, got 0x%016x", expected, actual)
		t.FailNow()
	}
}

func assertPanic(t *testing.T, f func(), errMsg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
			t.FailNow()
		} else if r != errMsg {
			t.Errorf(`Expected "%s", got "%s"`, errMsg, r)
			t.FailNow()
		}
	}()

	f()
}
