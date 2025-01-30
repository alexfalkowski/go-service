package encoding

import (
	"github.com/alexfalkowski/go-service/encoding/gob"
	"github.com/alexfalkowski/go-service/encoding/json"
	"github.com/alexfalkowski/go-service/encoding/proto"
	"github.com/alexfalkowski/go-service/encoding/toml"
	"github.com/alexfalkowski/go-service/encoding/yaml"
)

// Map of encoding.
type Map struct {
	encoders map[string]Encoder
}

// NewMap for encoding.
func NewMap() *Map {
	return &Map{
		encoders: map[string]Encoder{
			"json":      json.NewEncoder(),
			"yaml":      yaml.NewEncoder(),
			"yml":       yaml.NewEncoder(),
			"toml":      toml.NewEncoder(),
			"proto":     proto.NewBinary(),
			"protobuf":  proto.NewBinary(),
			"prototext": proto.NewText(),
			"protojson": proto.NewJSON(),
			"gob":       gob.NewEncoder(),
		},
	}
}

// Register kind and encoder.
func (f *Map) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get from kind.
func (f *Map) Get(kind string) Encoder {
	return f.encoders[kind]
}
