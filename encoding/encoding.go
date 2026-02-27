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

// MapParams defines the dependencies used to construct an encoding Map.
//
// It is intended for dependency injection (Fx/Dig). The default wiring is provided by `encoding.Module`.
type MapParams struct {
	di.In

	// JSON is the JSON encoder implementation registered under kind "json".
	JSON *json.Encoder

	// YAML is the YAML encoder implementation registered under kinds "yaml" and "yml".
	YAML *yaml.Encoder

	// TOML is the TOML encoder implementation registered under kind "toml".
	TOML *toml.Encoder

	// ProtoBinary is the protobuf binary encoder implementation registered under common binary kinds
	// (e.g. "proto", "protobuf", "pb", etc.).
	ProtoBinary *proto.Binary

	// ProtoText is the protobuf text encoder implementation registered under common text kinds
	// (e.g. "prototext", "prototxt", "pbtxt").
	ProtoText *proto.Text

	// ProtoJSON is the protobuf JSON encoder implementation registered under common JSON kinds
	// (e.g. "protojson", "pbjson").
	ProtoJSON *proto.JSON

	// GOB is the gob encoder implementation registered under kind "gob".
	GOB *gob.Encoder

	// Bytes is the passthrough encoder for io.ReaderFrom/io.WriterTo payloads, registered under kinds
	// like "plain", "octet-stream", and "markdown".
	Bytes *bytes.Encoder
}

// NewMap constructs a Map pre-populated with default encoders.
//
// The returned registry includes common kinds used throughout go-service, including:
//
//   - Structured config formats: "json", "yaml", "yml", "toml"
//
//   - Protobuf formats:
//
//   - binary: "proto", "protobuf", "pb", "protobin", "pbbin"
//
//   - text: "prototext", "prototxt", "pbtxt"
//
//   - JSON: "protojson", "pbjson"
//
//   - gob: "gob"
//
//   - bytes/plain passthrough: "plain", "octet-stream", "markdown"
//
// Callers can add additional kinds or override existing kinds via (*Map).Register.
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
//
// This type is a thin convenience around a string-keyed map and is commonly used with configuration
// to select an encoder at runtime.
//
// Map is not concurrency-safe. If you mutate it via Register, do so during initialization.
type Map struct {
	encoders map[string]Encoder
}

// Register associates kind with enc, overwriting any existing encoder.
//
// If kind already exists, the previous encoder is replaced.
func (f *Map) Register(kind string, enc Encoder) {
	f.encoders[kind] = enc
}

// Get returns the encoder registered for kind.
//
// If no encoder is registered for kind, Get returns nil. Callers typically treat nil as "unknown kind"
// and fall back to a default encoder elsewhere.
func (f *Map) Get(kind string) Encoder {
	return f.encoders[kind]
}

// Keys returns the list of registered encoder kinds.
//
// The returned slice is not guaranteed to be sorted.
func (f *Map) Keys() []string {
	return slices.Collect(maps.Keys(f.encoders))
}
