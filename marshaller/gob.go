package marshaller

import (
	"bytes"
	"encoding/gob"
)

// GOB for marshaller.
type GOB struct{}

// NewJSON for marshaller.
func NewGOB() *GOB {
	return &GOB{}
}

func (m *GOB) Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(v)

	return b.Bytes(), err
}

func (m *GOB) Unmarshal(data []byte, v any) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
