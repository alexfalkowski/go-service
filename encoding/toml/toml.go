package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// Marshaller for toml.
type Marshaller struct{}

// NewTOML for toml.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	err := toml.NewEncoder(&b).Encode(v)

	return b.Bytes(), err
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
