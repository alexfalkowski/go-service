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

	if err := gob.NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (m *GOB) Unmarshal(data []byte, v any) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
