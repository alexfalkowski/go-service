package toml

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/io"
)

var defaultEncoder = &Encoder{}

// NewEncoder constructs a TOML encoder.
//
// This encoder is a thin adapter around [github.com/BurntSushi/toml] that satisfies
// [github.com/alexfalkowski/go-service/v2/encoding.Encoder].
func NewEncoder() *Encoder {
	return defaultEncoder
}

// Encoder implements TOML encoding and decoding.
//
// It uses BurntSushi/toml and rejects undecoded keys.
type Encoder struct{}

// Encode writes v to w as TOML.
//
// This is a thin wrapper around `toml.NewEncoder(w).Encode(v)`.
func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

// Decode reads TOML from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
//
// This method rejects keys that do not decode into v.
func (e *Encoder) Decode(r io.Reader, v any) error {
	meta, err := toml.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		return fmt.Errorf("toml: undecoded key %s", undecoded[0])
	}

	return nil
}

// Marshal encodes v as TOML.
func Marshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := defaultEncoder.Encode(&buffer, v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Unmarshal decodes TOML data into v.
func Unmarshal(data []byte, v any) error {
	return defaultEncoder.Decode(bytes.NewReader(data), v)
}
