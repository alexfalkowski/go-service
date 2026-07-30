package stream

import (
	"maps"

	"github.com/alexfalkowski/go-service/v2/encoding/stream/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/yaml"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/slices"
)

// NewMap constructs a Map with the default streaming encoder/decoder constructors.
//
// The returned registry includes these kinds: "json", "msgpack", "gob", "yaml".
//
// Each per-kind package (encoding/stream/json, .../msgpack, .../gob, .../yaml) exports constructors
// returning its own concrete [json.Encoder]/[json.Decoder]-shaped type rather than the
// [Encoder]/[Decoder] interfaces, so those packages do not import this one back — avoiding a
// compile-time import cycle, since this package imports them here to build the default registry. The
// adapter closures below convert each concrete constructor to the registry's constructor-function
// shape; the concrete types satisfy [Encoder]/[Decoder] structurally, so the conversion is just an
// ordinary interface assignment on return.
//
// Callers can add additional kinds or override existing kinds via [Map.RegisterEncoder] and
// [Map.RegisterDecoder].
func NewMap() *Map {
	return &Map{
		encoders: map[string]EncoderFunc{
			"json":    func(w io.Writer) Encoder { return json.NewEncoder(w) },
			"msgpack": func(w io.Writer) Encoder { return msgpack.NewEncoder(w) },
			"gob":     func(w io.Writer) Encoder { return gob.NewEncoder(w) },
			"yaml":    func(w io.Writer) Encoder { return yaml.NewEncoder(w) },
		},
		decoders: map[string]DecoderFunc{
			"json":    func(r io.Reader) Decoder { return json.NewDecoder(r) },
			"msgpack": func(r io.Reader) Decoder { return msgpack.NewDecoder(r) },
			"gob":     func(r io.Reader) Decoder { return gob.NewDecoder(r) },
			"yaml":    func(r io.Reader) Decoder { return yaml.NewDecoder(r) },
		},
	}
}

// Map provides lookup and registration of streaming encoder/decoder constructors by kind.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding.Map], which registers ready-to-use
// encoder instances, Map registers constructor functions. Each streaming Encode/Decode call needs a
// fresh [Encoder] or [Decoder] bound to its own writer or reader, so the registry stores how to build
// one rather than a shared instance.
//
// Map is not concurrency-safe. If you mutate it via RegisterEncoder or RegisterDecoder, do so during
// initialization.
type Map struct {
	encoders map[string]EncoderFunc
	decoders map[string]DecoderFunc
}

// RegisterEncoder associates kind with fn, overwriting any existing encoder constructor.
//
// If kind already exists, the previous encoder constructor is replaced.
func (m *Map) RegisterEncoder(kind string, fn EncoderFunc) {
	m.encoders[kind] = fn
}

// RegisterDecoder associates kind with fn, overwriting any existing decoder constructor.
//
// If kind already exists, the previous decoder constructor is replaced.
func (m *Map) RegisterDecoder(kind string, fn DecoderFunc) {
	m.decoders[kind] = fn
}

// GetEncoder returns the encoder constructor registered for kind.
//
// If no encoder constructor is registered for kind, or if kind was registered with a nil
// constructor, GetEncoder returns nil. Callers typically treat nil as "unknown or unavailable kind"
// and fail explicitly rather than falling back to a default encoder.
//
// GetEncoder is nil-safe: it returns nil for a nil m, so resolving against an unconfigured registry
// fails the same way resolving an unregistered kind does, rather than panicking.
func (m *Map) GetEncoder(kind string) EncoderFunc {
	if m == nil {
		return nil
	}

	return m.encoders[kind]
}

// GetDecoder returns the decoder constructor registered for kind.
//
// If no decoder constructor is registered for kind, or if kind was registered with a nil
// constructor, GetDecoder returns nil. Callers typically treat nil as "unknown or unavailable kind"
// and fail explicitly rather than falling back to a default decoder.
//
// GetDecoder is nil-safe: it returns nil for a nil m, so resolving against an unconfigured registry
// fails the same way resolving an unregistered kind does, rather than panicking.
func (m *Map) GetDecoder(kind string) DecoderFunc {
	if m == nil {
		return nil
	}

	return m.decoders[kind]
}

// Keys returns the union of registered encoder and decoder kinds.
//
// Keys includes kinds registered with nil constructors. The returned slice is not guaranteed to be
// sorted.
func (m *Map) Keys() []string {
	keys := make(map[string]struct{}, len(m.encoders)+len(m.decoders))
	for k := range m.encoders {
		keys[k] = struct{}{}
	}

	for k := range m.decoders {
		keys[k] = struct{}{}
	}

	return slices.Collect(maps.Keys(keys))
}
