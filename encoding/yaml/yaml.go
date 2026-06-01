package yaml

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/errors"
	"github.com/alexfalkowski/go-service/v2/io"
	yaml "go.yaml.in/yaml/v3"
)

var defaultEncoder = &Encoder{}

// NewEncoder constructs a YAML encoder.
//
// This encoder is a thin adapter around go-yaml v3 (imported as go.yaml.in/yaml/v3) that satisfies
// `github.com/alexfalkowski/go-service/v2/encoding.Encoder`.
func NewEncoder() *Encoder {
	return defaultEncoder
}

// Encoder implements YAML encoding and decoding.
//
// It uses go-yaml v3 with default settings.
type Encoder struct{}

// Encode writes v to w as YAML.
//
// Encode closes the upstream YAML encoder after writing so final buffered data and
// finalization errors are handled.
func (e *Encoder) Encode(w io.Writer, v any) error {
	encoder := yaml.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		return err
	}

	return encoder.Close()
}

// Decode reads YAML from r and decodes it into v.
//
// In most cases v should be a pointer to the destination value (for example *MyStruct).
// Decode reads one YAML document and rejects additional documents in the same stream.
func (e *Encoder) Decode(r io.Reader, v any) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	var extra any
	return errors.TrailingData(decoder.Decode(&extra))
}

// Marshal encodes v as YAML.
func Marshal(v any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := defaultEncoder.Encode(&buffer, v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Unmarshal decodes one YAML document from data into v.
//
// It uses Decode, so it rejects additional YAML documents after the first payload.
func Unmarshal(data []byte, v any) error {
	return defaultEncoder.Decode(bytes.NewReader(data), v)
}
