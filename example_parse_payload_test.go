package bits_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
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
func (r *PayloadReader) ReadTime() (time.Time, error) {
	if r.Err != nil {
		return time.Time{}, r.Err
	}

	if b, err := r.ReadBits(32); err != nil {
		r.Err = errors.WithMessage(err, "read time")
		return time.Time{}, r.Err
	} else {
		return time.Unix(int64(b), 0).UTC(), nil
	}
}

// ParsePayload creates a Payload from h, a hex-encoded string.
func ParsePayload(h string) (p *Payload, err error) {
	var d []byte
	d, err = hex.DecodeString(h)
	if err != nil {
		return
	}

	var r = NewPayloadReader(d)

	// This block of code directly describes the format of the payload.
	p = &Payload{}

	// The following are examples of two styles of using a Reader to parse a
	// binary payload. In the first, errors are checked after each read. This
	// allows the implementation to immediately address a problem. In the
	// second, checking for errors is deferred until the end, allowing for a
	// more concise implementation.

	// Handling errors immediately.
	var b bits.Block
	if b, err = r.ReadBits(4); err != nil {
		return
	} else {
		p.Version = int(b)
	}
	if b, err = r.ReadBits(8); err != nil {
		return
	} else {
		p.Category = Category(b)
	}
	if p.IsX, err = r.ReadBool(); err != nil {
		return
	}

	// Deferring error checks.
	p.IsY, _ = r.ReadBool()
	p.IsZ, _ = r.ReadBool()
	p.Created, _ = r.ReadTime()
	p.Modified, _ = r.ReadTime()

	return p, r.Err
}

func Example_parsePayload() {
	// In this example, the payload is presented as a hex-encoded string.
	var p, err = ParsePayload("b03a877346aa877346aa")

	// Pretty print in JSON for readability.
	var pretty, _ = json.MarshalIndent(p, "", "  ")
	fmt.Println(string(pretty))
	fmt.Println(err)

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
	// <nil>
}

func Example_parsePayloadError() {
	// The payload has been truncated and does not contain enough bits.
	var p, err = ParsePayload("b03a8773")

	// Pretty print in JSON for readability.
	var pretty, _ = json.MarshalIndent(p, "", "  ")
	fmt.Println(string(pretty))
	fmt.Println(err)

	// Output:
	// {
	//   "Version": 11,
	//   "Category": 3,
	//   "IsX": true,
	//   "IsY": false,
	//   "IsZ": true,
	//   "Created": "0001-01-01T00:00:00Z",
	//   "Modified": "0001-01-01T00:00:00Z"
	// }
	// read time: read bits (index=15, length=32): bits: length extends beyond range
}
