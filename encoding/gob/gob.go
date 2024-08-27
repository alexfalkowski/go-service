package gob

import (
	"bytes"
	"encoding/gob"
	"io"
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

// Encoder for gob.
type Encoder struct{}

// NewEncoder for gob.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return gob.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}
