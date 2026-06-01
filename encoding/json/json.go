package json

import (
	"encoding/json"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/strings"
)

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

var defaultEncoder = &Encoder{}

// NewEncoder constructs a JSON encoder.
//
// NewEncoder returns an [Encoder] that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder` while delegating to
// the standard library's `encoding/json` implementation with readable indented
// encoding and default decoding.
func NewEncoder() *Encoder {
	return defaultEncoder
}

// Encoder implements JSON encoding and decoding.
//
// It writes readable indented JSON and uses the standard library
// `encoding/json` decoder.
type Encoder struct{}

// Encode writes v to w as JSON.
//
// Encode is equivalent to calling `json.NewEncoder(w).SetIndent("", "  ")`
// before encoding v. As with the standard library, the encoded output is
// terminated with a trailing newline.
func (e *Encoder) Encode(w io.Writer, v any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent(strings.Empty, "  ")

	return encoder.Encode(v)
}

// Decode reads JSON from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value, such as
// `*MyStruct` or `*map[string]any`. Decode is equivalent to calling
// `json.NewDecoder(r).Decode(v)`, then requiring the stream to contain no
// additional JSON values.
//
// Duplicate JSON object keys keep the standard library's last-wins behavior.
func (e *Encoder) Decode(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	var extra any
	return errors.TrailingData(decoder.Decode(&extra))
}

// Marshal encodes v as readable indented JSON.
//
// It exists so repository packages can use a single go-service JSON import path.
func Marshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, strings.Empty, "  ")
}

// Unmarshal decodes one JSON value from data into v.
//
// It uses Decode, so it rejects additional JSON values after the first payload.
func Unmarshal(data []byte, v any) error {
	return defaultEncoder.Decode(bytes.NewReader(data), v)
}

// Valid reports whether data is syntactically valid JSON.
//
// It behaves exactly like `encoding/json.Valid`.
func Valid(data []byte) bool {
	return json.Valid(data)
}
