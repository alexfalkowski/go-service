package stream

import (
	"maps"

	"github.com/alexfalkowski/go-service/v2/encoding/codec"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/json"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/stream/yaml"
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/slices"
)

// NewMap constructs a Map with the default streaming codecs.
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
// Callers can add additional kinds or override existing kinds via [Map.Register].
func NewMap() *Map {
	return &Map{
		codecs: map[string]Codec{
			"json": {
				Encoder: func(w io.Writer) Encoder { return json.NewEncoder(w) },
				Decoder: func(r io.Reader, opts ...codec.Option) Decoder { return json.NewDecoder(r, opts...) },
			},
			"msgpack": {
				Encoder: func(w io.Writer) Encoder { return msgpack.NewEncoder(w) },
				Decoder: func(r io.Reader, opts ...codec.Option) Decoder { return msgpack.NewDecoder(r, opts...) },
			},
			"gob": {
				Encoder: func(w io.Writer) Encoder { return gob.NewEncoder(w) },
				Decoder: func(r io.Reader, opts ...codec.Option) Decoder { return gob.NewDecoder(r, opts...) },
			},
			"yaml": {
				Encoder: func(w io.Writer) Encoder { return yaml.NewEncoder(w) },
				Decoder: func(r io.Reader, opts ...codec.Option) Decoder { return yaml.NewDecoder(r, opts...) },
			},
		},
	}
}

// Codec bundles the encoder and decoder constructors registered for one kind.
//
// Either field may be nil when the kind supports only one direction; callers must check the field they
// intend to use before calling it.
type Codec struct {
	// Encoder constructs an [Encoder] bound to a writer, or nil if this kind has no registered encoder.
	Encoder EncoderFunc

	// Decoder constructs a [Decoder] bound to a reader, or nil if this kind has no registered decoder.
	Decoder DecoderFunc
}

// Map provides lookup and registration of streaming codecs by kind.
//
// Unlike [github.com/alexfalkowski/go-service/v2/encoding.Map], which registers ready-to-use
// encoder instances, Map registers constructor functions. Each streaming Encode/Decode call needs a
// fresh [Encoder] or [Decoder] bound to its own writer or reader, so the registry stores how to build
// one rather than a shared instance.
//
// Map is not concurrency-safe. If you mutate it via Register, do so during initialization.
type Map struct {
	codecs map[string]Codec
}

// Register associates kind with codec, overwriting any existing codec in full.
//
// If kind already exists, the previous codec is replaced entirely — Register does not merge codec's
// Encoder/Decoder fields with whatever was previously registered, so registering a partial Codec (for
// example, only Encoder set) clears the other field for kind.
func (m *Map) Register(kind string, codec Codec) {
	m.codecs[kind] = codec
}

// Get returns the codec registered for kind.
//
// If no codec is registered for kind, Get returns the zero [Codec] (both fields nil). Callers typically
// treat a nil Encoder or Decoder field as "unknown or unavailable kind" and fail explicitly rather than
// falling back to a default.
//
// Get is nil-safe: it returns the zero Codec for a nil m, so resolving against an unconfigured registry
// fails the same way resolving an unregistered kind does, rather than panicking.
func (m *Map) Get(kind string) Codec {
	if m == nil {
		return Codec{}
	}

	return m.codecs[kind]
}

// Keys returns the registered kinds.
//
// Keys includes kinds registered with a zero or partially-nil [Codec]. The returned slice is not
// guaranteed to be sorted.
func (m *Map) Keys() []string {
	return slices.Collect(maps.Keys(m.codecs))
}
