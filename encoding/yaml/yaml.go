package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

// Marshaller for yaml.
type Marshaller struct{}

// NewMarshaller for yaml.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

// Encoder for yaml.
type Encoder struct{}

// NewEncoder for yaml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return yaml.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return yaml.NewDecoder(r).Decode(v)
}
