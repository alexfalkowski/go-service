package hjson

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/io"
	hjson "github.com/hjson/hjson-go/v4"
)

var defaultEncoder = &Encoder{}

// NewEncoder constructs an HJSON encoder.
//
// This encoder is a thin adapter around [github.com/hjson/hjson-go] that satisfies
// [github.com/alexfalkowski/go-service/v2/encoding/codec.Encoder].
func NewEncoder() *Encoder {
	return defaultEncoder
}

// Encoder implements HJSON encoding and decoding.
type Encoder struct{}

// Encode writes v to w as HJSON. HJSON has no required-field concept, so
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] has no effect.
func (e *Encoder) Encode(w io.Writer, v any, _ ...codec.Option) error {
	data, err := hjson.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// Decode reads HJSON from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// Decode rejects duplicate object keys and unknown destination fields unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] is
// supplied.
// Decode buffers all remaining input from r before decoding; callers should bound
// untrusted or potentially large readers before calling it.
func (e *Encoder) Decode(r io.Reader, v any, opts ...codec.Option) error {
	data, _, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	resolved := codec.Apply(opts...)

	options := hjson.DefaultDecoderOptions()
	options.DisallowDuplicateKeys = true
	options.DisallowUnknownFields = !resolved.DiscardUnknown()

	return hjson.UnmarshalWithOptions(data, v, options)
}

// Marshal encodes v as HJSON.
func Marshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := defaultEncoder.Encode(&buffer, v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Unmarshal decodes HJSON data into v.
//
// It uses Decode, so it rejects duplicate object keys and unknown destination fields.
func Unmarshal(data []byte, v any) error {
	return defaultEncoder.Decode(bytes.NewReader(data), v)
}
