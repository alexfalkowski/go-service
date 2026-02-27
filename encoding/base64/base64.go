package base64

import (
	"encoding/base64"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Encode encodes src using standard base64 encoding (RFC 4648) and returns the encoded string.
//
// This helper uses `base64.StdEncoding` and allocates a new buffer sized to the encoded output.
func Encode(src []byte) string {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)

	return bytes.String(buf)
}

// Decode decodes a standard base64-encoded string s into a byte slice.
//
// This helper uses `base64.StdEncoding`. It allocates a destination buffer sized to the maximum decoded
// length, then returns the subslice actually written.
//
// It returns a non-nil error if s contains invalid base64 data.
func Decode(s string) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(buf, strings.Bytes(s))

	return buf[:n], err
}
