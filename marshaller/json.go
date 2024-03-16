package marshaller

import (
	"encoding/json"
)

// JSON for marshaller.
type JSON struct{}

// NewJSON for marshaller.
func NewJSON() *JSON {
	return &JSON{}
}

func (m *JSON) Marshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "    ")
}

func (m *JSON) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
