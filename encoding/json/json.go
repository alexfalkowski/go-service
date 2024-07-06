package json

import (
	"github.com/goccy/go-json"
)

// Marshaller for json.
type Marshaller struct{}

// NewMarshaller for json.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
