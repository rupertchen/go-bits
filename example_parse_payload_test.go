package bits_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rupertchen/go-bits"
)

type Category int8

// Payload defines the properties of a payload.
type Payload struct {
	Version       int       // 4 bits
	Category      Category  // 8 bits
	IsX, IsY, IsZ bool      // 1 bit each
	Created       time.Time // 32 bits, seconds since epoch
	Modified      time.Time // 32 bits, seconds since epoch
}

// PayloadReader provides convenient typed Read* methods for parsing a Payload
// from raw bits.
type PayloadReader struct {
	*bits.Reader
}

// NewPayloadReader returns a PayloadReader that is ready to parse a Payload
// from src.
func NewPayloadReader(src []byte) *PayloadReader {
	return &PayloadReader{bits.NewReader(bits.NewBitmap(src))}
}

// ReadTime interprets the next 32 bits as the Unix epoch and returns a time.
// Time in UTC.
func (r *PayloadReader) ReadTime() time.Time {
	var sec = int64(r.ReadBits(32))
	return time.Unix(sec, 0).UTC()
}

// ParsePayload creates a Payload from h, a hex-encoded string.
func ParsePayload(h string) (p *Payload, err error) {
	defer func() {
		if r := recover(); r!= nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	var b []byte
	b, err = hex.DecodeString(h)
	if err != nil {
		return
	}

	var r = NewPayloadReader(b)

	// This block of code directly describes the format of the payload.
	p = &Payload{}
	p.Version = int(r.ReadBits(4))
	p.Category = Category(r.ReadBits(8))
	p.IsX = r.ReadBool()
	p.IsY = r.ReadBool()
	p.IsZ = r.ReadBool()
	p.Created = r.ReadTime()
	p.Modified = r.ReadTime()

	return p, nil
}

func Example_parsePayload() {
	// In this example, the payload is presented as a hex-encoded string.
	var p, _ = ParsePayload("b03a877346aa877346aa")

	// Pretty print in JSON for readability.
	var pretty, _ = json.MarshalIndent(p, "", "  ")
	fmt.Println(string(pretty))

	// Output:
	// {
	//   "Version": 11,
	//   "Category": 3,
	//   "IsX": true,
	//   "IsY": false,
	//   "IsZ": true,
	//   "Created": "2006-01-02T22:04:05Z",
	//   "Modified": "2006-01-02T22:04:05Z"
	// }
}
