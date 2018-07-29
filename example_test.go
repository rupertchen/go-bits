package bits_test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/rupertchen/go-bits"
)

func ExampleNewBitmap_byteSlice() {
	// Load bytes from a []byte.
	var bmp = bits.NewBitmap([]byte("Hello, World!"))

	fmt.Printf("%c%c%c%c%c%c\n",
		bmp.MustGet(0, 8),
		bmp.MustGet(8, 8),
		bmp.MustGet(16, 8),
		bmp.MustGet(24, 8),
		bmp.MustGet(32, 8),
		bmp.MustGet(96, 8),
	)

	// Output: Hello!
}

func ExampleNewBitmap_base64() {
	// Load bytes from a base64 encoded string.
	var b64 = "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQ="
	var decoded, err = base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatal(err)
	}
	var bmp = bits.NewBitmap(decoded)

	for i := 0; i < 11; i++ {
		fmt.Printf("%c", bmp.MustGet(uint(8*i), 8))
	}
	fmt.Print("\n")

	// Output: Lorem ipsum
}

func ExampleNewBitmap_hexString() {
	// Load bytes from a hex encoded string.
	var s = "a40c9a21e5a1"
	var decoded, err = hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	var bmp = bits.NewBitmap(decoded)

	fmt.Printf("0x%03x\n", bmp.MustGet(0, 12))
	fmt.Printf("0x%03x\n", bmp.MustGet(12, 12))
	fmt.Printf("0x%03x\n", bmp.MustGet(24, 12))
	fmt.Printf("0x%03x\n", bmp.MustGet(36, 12))

	// Output:
	// 0xa40
	// 0xc9a
	// 0x21e
	// 0x5a1
}

func ExampleNewBitmapFromBlocks() {
	var bmp = bits.NewBitmapFromBlocks([]bits.Block{
		0x0123456789abcdef,
		0xfedcba9876543210,
	})

	fmt.Printf("0x%08x\n", bmp.MustGet(0, 32))
	fmt.Printf("0x%08x\n", bmp.MustGet(32, 32))
	fmt.Printf("0x%08x\n", bmp.MustGet(64, 32))
	fmt.Printf("0x%08x\n", bmp.MustGet(96, 32))

	// Output:
	// 0x01234567
	// 0x89abcdef
	// 0xfedcba98
	// 0x76543210
}

func ExampleBitmap_MustGet() {
	var bmp = bits.NewBitmap([]byte{0xab, 0xcd, 0xef})

	fmt.Printf("0x%016x\n", bmp.MustGet(0, 0))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 4))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 8))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 12))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 16))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 20))
	fmt.Printf("0x%016x\n", bmp.MustGet(0, 24))

	// Output:
	// 0x0000000000000000
	// 0x000000000000000a
	// 0x00000000000000ab
	// 0x0000000000000abc
	// 0x000000000000abcd
	// 0x00000000000abcde
	// 0x0000000000abcdef
}
