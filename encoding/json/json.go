package json

import (
	"encoding/json"
	"io"
)

// NewEncoder constructs a JSON encoder.
//
// NewEncoder returns an [Encoder] that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder` while delegating to
// the standard library's `encoding/json` implementation with its default
// settings.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Marshaler aliases the standard library JSON marshaler interface.
//
// Use this alias when a package wants to refer to the marshaling contract while
// keeping imports within the go-service JSON package.
type Marshaler = json.Marshaler

// Unmarshaler aliases the standard library JSON unmarshaler interface.
//
// Use this alias when a package wants to refer to the unmarshaling contract
// while keeping imports within the go-service JSON package.
type Unmarshaler = json.Unmarshaler

// RawMessage aliases the standard library raw JSON message type.
//
// It is useful when callers need to defer decoding or preserve the original
// encoded JSON payload.
type RawMessage = json.RawMessage

// Number aliases the standard library JSON number type.
//
// It is primarily useful with decoders configured to preserve numeric values as
// strings until the caller decides how to interpret them.
type Number = json.Number

// Encoder implements JSON encoding and decoding.
//
// It preserves the default behavior of the standard library `encoding/json`
// encoder and decoder and does not add repository-specific options.
type Encoder struct{}

// Encode writes v to w as JSON.
//
// Encode is equivalent to calling `json.NewEncoder(w).Encode(v)`. As with the
// standard library, the encoded output is terminated with a trailing newline.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

// Decode reads JSON from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value, such as
// `*MyStruct` or `*map[string]any`. Decode is equivalent to calling
// `json.NewDecoder(r).Decode(v)`.
func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

// Marshal encodes v as JSON using the standard library implementation.
//
// It exists so repository packages can use a single go-service JSON import path
// without changing standard library marshaling behavior.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent encodes v as indented JSON using the standard library
// implementation.
//
// It behaves exactly like `encoding/json.MarshalIndent`.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal decodes JSON data into v using the standard library implementation.
//
// In most cases v should be a pointer to the destination value. The decoding
// rules and error behavior are identical to `encoding/json.Unmarshal`.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Valid reports whether data is syntactically valid JSON.
//
// It behaves exactly like `encoding/json.Valid`.
func Valid(data []byte) bool {
	return json.Valid(data)
}
