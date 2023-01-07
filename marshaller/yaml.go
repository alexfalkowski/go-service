package marshaller

import (
	"gopkg.in/yaml.v3"
)

// YAML for marshaller.
type YAML struct{}

// NewYAML for marshaller.
func NewYAML() *YAML {
	return &YAML{}
}

func (m *YAML) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (m *YAML) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
