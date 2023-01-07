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

	if err := toml.NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (m *TOML) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
