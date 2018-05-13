package bits

import "testing"

func TestScanner_ReadBits(t *testing.T) {
	var bmp = NewBitmap(128)
	bmp.store[0] = 0xaaaaaaaaaaaaaaaa
	bmp.store[1] = 0x5555555555555555
	var s = NewScanner(bmp)

	boolEquals(t, true, s.ReadBool())
	boolEquals(t, false, s.ReadBool())
	boolEquals(t, true, s.ReadBool())
	boolEquals(t, false, s.ReadBool())
	byteEquals(t, 0xaa, s.ReadByte())
	boolEquals(t, true, s.ReadBool())
	byteEquals(t, 0x55, s.ReadByte())
	blockEquals(t, 0x0000002aaaaaaaaa, s.ReadBits(39))
	byteEquals(t, 0xa5, s.ReadByte())
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
