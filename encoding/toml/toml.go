package toml

import (
	"io"

	"github.com/BurntSushi/toml"
)

// Marshaller for toml.
type Marshaller struct{}

// NewTOML for toml.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}

// Encoder for toml.
type Encoder struct{}

// NewEncoder for toml.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return toml.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	_, err := toml.NewDecoder(r).Decode(v)

	return err
}
