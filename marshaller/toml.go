package marshaller

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// TOML for marshaller.
type TOML struct{}

// NewTOML for marshaller.
func NewTOML() *TOML {
	return &TOML{}
}

func (m *TOML) Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	err := toml.NewEncoder(&b).Encode(v)

	return b.Bytes(), err
}

func (m *TOML) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
