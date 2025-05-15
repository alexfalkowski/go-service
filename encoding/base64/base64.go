package base64

import "encoding/base64"

// Encode src to a string.
func Encode(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Decode str to []byte, otherwise err if fails.
func Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
