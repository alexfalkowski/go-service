package yaml

import (
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
