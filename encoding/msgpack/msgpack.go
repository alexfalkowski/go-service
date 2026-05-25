package msgpack

import (
	"github.com/Basekick-Labs/msgpack/v6"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
)

// NewEncoder constructs a MessagePack encoder.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encoder implements MessagePack encoding and decoding.
type Encoder struct{}

// Encode writes v to w as MessagePack.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return msgpack.NewEncoder(w).Encode(v)
}

// Decode reads MessagePack from r and decodes it into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	decoder := msgpack.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	_, err := decoder.PeekCode()
	return errors.TrailingData(err)
}

// Marshal encodes v as MessagePack.
func Marshal(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

// Unmarshal decodes MessagePack data into v.
func Unmarshal(data []byte, v any) error {
	return NewEncoder().Decode(bytes.NewReader(data), v)
}
