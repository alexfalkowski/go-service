package json

import (
	"io"

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

// Encoder for json.
type Encoder struct{}

// NewEncoder for json.
func NewEncoder() *Encoder {
	return &Encoder{}
}

func (e *Encoder) Encode(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (e *Encoder) Decode(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}
