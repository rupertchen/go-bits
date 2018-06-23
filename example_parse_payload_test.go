package bits_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rupertchen/go-bits"
	"github.com/pkg/errors"
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
func (r *PayloadReader) ReadTime() (time.Time, error) {
	if (r.Err != nil) {
		return time.Time{}, r.Err
	}

	if b, err := r.ReadBits(32); err != nil {
		r.Err = errors.Wrap(err, "read time")
		return time.Time{}, r.Err
	} else {
		return time.Unix(int64(b), 0).UTC(), nil
	}
}

// ParsePayload creates a Payload from h, a hex-encoded string.
func ParsePayload(h string) (p *Payload, err error) {
	var b []byte
	b, err = hex.DecodeString(h)
	if err != nil {
		return
	}

	var r = NewPayloadReader(b)

	// This block of code directly describes the format of the payload.
	p = &Payload{}
	if b, err := r.ReadBits(4); err == nil {
		p.Version = int(b)
	}
	if b, err := r.ReadBits(8); err == nil {
		p.Category = Category(b)
	}
	p.IsX, err = r.ReadBool()
	p.IsY, err = r.ReadBool()
	p.IsZ, err = r.ReadBool()
	p.Created, err = r.ReadTime()
	p.Modified, err = r.ReadTime()

	return p, err
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
