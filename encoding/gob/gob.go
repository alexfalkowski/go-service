package gob

import (
	"bytes"
	"encoding/gob"
)

// Marshaller for gob.
type Marshaller struct{}

// NewMarshaller for gob.
func NewMarshaller() *Marshaller {
	return &Marshaller{}
}

func (m *Marshaller) Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)

	return b.Bytes(), err
}

func (m *Marshaller) Unmarshal(data []byte, v any) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
