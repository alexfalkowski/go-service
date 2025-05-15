package base64

import (
	"encoding/base64"

	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/strings"
)

// Encode src to a string.
func Encode(src []byte) string {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)

	return bytes.String(buf)
}

// Decode str to []byte, otherwise err if fails.
func Decode(s string) ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(buf, strings.Bytes(s))

	return buf[:n], err
}
