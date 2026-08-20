package msgpack

import (
	"github.com/Basekick-Labs/msgpack/v6"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

var defaultEncoder = &Encoder{}

// NewEncoder constructs a MessagePack encoder.
func NewEncoder() *Encoder {
	return defaultEncoder
}

// Encoder implements MessagePack encoding and decoding.
type Encoder struct{}

// Encode writes v to w as MessagePack. MessagePack has no required-field concept, so
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithAllowPartial] has no effect.
func (e *Encoder) Encode(w io.Writer, v any, _ ...codec.Option) error {
	return msgpack.NewEncoder(w).Encode(v)
}

// Decode reads one MessagePack value from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// Decode rejects unknown destination fields unless
// [github.com/alexfalkowski/go-service/v2/encoding/codec.WithDiscardUnknown] is
// supplied. It also rejects trailing encoded values or malformed trailing data.
func (e *Encoder) Decode(r io.Reader, v any, opts ...codec.Option) error {
	decoder := msgpack.NewDecoder(r)
	resolved := codec.Apply(opts...)
	decoder.DisallowUnknownFields(!resolved.DiscardUnknown())

	if err := decoder.Decode(v); err != nil {
		return err
	}

	_, err := decoder.PeekCode()
	return errors.TrailingData(err)
}

// Marshal encodes v as MessagePack.
func Marshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := defaultEncoder.Encode(&buffer, v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Unmarshal decodes one MessagePack value from data into v.
//
// It uses Decode, so it rejects trailing encoded values or malformed trailing data.
func Unmarshal(data []byte, v any) error {
	return defaultEncoder.Decode(bytes.NewReader(data), v)
}
