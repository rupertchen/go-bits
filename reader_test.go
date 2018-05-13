package bits

import "testing"

func TestReader_ReadBits(t *testing.T) {
	var r = NewReader(NewBitmapFromBlocks([]Block{0xaaaaaaaaaaaaaaaa, 0x5555555555555555}))

	var tests = []struct {
		expected Block
		length   uint
	}{
		// Read one bit at a time to verify position changes.
		{0x1, 1},
		{0x0, 1},
		{0x1, 1},
		{0x0, 1},

		// Read multiple bits at arbitrary alignments.
		{0xaa, 8},
		{0x1, 1},
		{0x55, 8},

		// Read multiple bits of non-word and non-byte aligned length.
		{0x2aaaaaaaaa, 39},

		// Read spanning internal blocks.
		{0xa5, 8},
	}
	for _, test := range tests {
		blockEquals(t, test.expected, r.ReadBits(test.length))
	}
}

func TestReader_ReadBool(t *testing.T) {
	var r = NewReader(NewBitmap([]byte{0xa5}))

	var expecteds = []bool{
		true, false, true, false,
		false, true, false, true,
	}
	for _, expected := range expecteds {
		boolEquals(t, expected, r.ReadBool())
	}
}

func TestReader_ReadByte(t *testing.T) {
	var r = NewReader(NewBitmapFromBlocks([]Block{0xa5fe735d00000000}))

	var expecteds = []byte{0xa5, 0xfe, 0x73, 0x5d}
	for _, expected := range expecteds {
		byteEquals(t, expected, r.ReadByte())
	}
}

func boolEquals(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Errorf("Expected %t, got %t", expected, actual)
		t.FailNow()
	}
}

func byteEquals(t *testing.T, expected, actual byte) {
	if expected != actual {
		t.Errorf("Expected 0x%02x, got 0x%02x", expected, actual)
	}
}
