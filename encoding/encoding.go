package encoding

import (
	"maps"
	"slices"

	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
)

// MapParams defines dependencies used to construct an encoding Map.
type MapParams struct {
	di.In
	JSON        *json.Encoder
	YAML        *yaml.Encoder
	TOML        *toml.Encoder
	ProtoBinary *proto.Binary
	ProtoText   *proto.Text
	ProtoJSON   *proto.JSON
	GOB         *gob.Encoder
	Bytes       *bytes.Encoder
}

// NewMap constructs a Map pre-populated with default encoders.
func NewMap(params MapParams) *Map {
	return &Map{
		encoders: map[string]Encoder{
			"json":         params.JSON,
			"yaml":         params.YAML,
			"yml":          params.YAML,
			"toml":         params.TOML,
			"pb":           params.ProtoBinary,
			"pbbin":        params.ProtoBinary,
			"proto":        params.ProtoBinary,
			"protobin":     params.ProtoBinary,
			"protobuf":     params.ProtoBinary,
			"pbtxt":        params.ProtoText,
			"prototext":    params.ProtoText,
			"prototxt":     params.ProtoText,
			"protojson":    params.ProtoJSON,
			"pbjson":       params.ProtoJSON,
			"gob":          params.GOB,
			"markdown":     params.Bytes,
			"octet-stream": params.Bytes,
			"plain":        params.Bytes,
		},
	}
}

// Map provides lookup and registration of encoders by kind.
type Map struct {
	encoders map[string]Encoder
}

// Register associates kind with enc, overwriting any existing encoder.
func (f *Map) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get returns the encoder registered for kind.
func (f *Map) Get(kind string) Encoder {
	return f.encoders[kind]
}

// Keys returns the list of registered encoder kinds.
func (f *Map) Keys() []string {
	return slices.Collect(maps.Keys(f.encoders))
}
