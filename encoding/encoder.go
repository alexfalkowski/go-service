package encoding

import (
	"io"

	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
)

// Encoder allows different types of encoding/decoding.
type Encoder interface {
	// Encode any to a writer.
	Encode(w io.Writer, e any) error

	// Decode any from a reader.
	Decode(r io.Reader, e any) error
}

type encoders map[string]Encoder

// MarshallerMap of marshaller.
type EncoderMap struct {
	encoders encoders
}

// NewMarshallerMap for marshaller.
func NewEncoderMap() *EncoderMap {
	m := &EncoderMap{
		encoders: encoders{
			"json":     json.NewEncoder(),
			"yaml":     yaml.NewEncoder(),
			"yml":      yaml.NewEncoder(),
			"toml":     toml.NewEncoder(),
			"proto":    proto.NewEncoder(),
			"protobuf": proto.NewEncoder(),
			"gob":      gob.NewEncoder(),
		},
	}

	return m
}

// Register kind and encoder.
func (f *EncoderMap) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get from kind.
func (f *EncoderMap) Get(kind string) Encoder {
	return f.encoders[kind]
}
