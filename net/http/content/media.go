package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/mime"
	content "github.com/elnormous/contenttype"
)

const (
	jsonKind     = "json"
	errorSubtype = "error"
)

// NewMedia builds a Media from a parsed media type and encoder map.
//
// Encoder selection:
//   - If the subtype is "error", it returns a Media without an encoder.
//   - If no encoder is registered for the subtype, it falls back to JSON.
//   - Otherwise it returns the encoder registered for the subtype.
func NewMedia(media content.MediaType, enc *encoding.Map) Media {
	return newMedia(media, enc)
}

// Media describes an HTTP media type and its associated encoder.
//
// Type is the full media type string (including parameters if present). Subtype is the parsed media subtype used
// for encoder lookup. Encoder may be nil when Subtype is "error".
type Media struct {
	// Encoder is the encoder/decoder associated with the media subtype.
	Encoder encoding.Encoder

	// Type is the full media type string (for example "application/json", "application/hjson", or "text/plain; charset=utf-8").
	Type string

	// Subtype is the parsed subtype (for example "json", "hjson", or "plain").
	Subtype string
}

// IsError reports whether the media subtype represents an error payload.
func (t Media) IsError() bool {
	return t.Subtype == errorSubtype
}

func knownMedia(mediaType string, enc *encoding.Map) (Media, bool) {
	// Exact built-in media types avoid the general parser on hot request paths.
	// Parameterized values still use the parser so their normalized Type string stays unchanged.
	switch mediaType {
	case mime.JSONMediaType:
		return jsonMedia(enc), true
	case mime.HJSONMediaType:
		return Media{Type: mime.HJSONMediaType, Subtype: "hjson", Encoder: enc.Get("hjson")}, true
	case mime.YAMLMediaType:
		return Media{Type: mime.YAMLMediaType, Subtype: "yaml", Encoder: enc.Get("yaml")}, true
	case mime.TOMLMediaType:
		return Media{Type: mime.TOMLMediaType, Subtype: "toml", Encoder: enc.Get("toml")}, true
	case mime.ProtobufMediaType:
		return Media{Type: mime.ProtobufMediaType, Subtype: "protobuf", Encoder: enc.Get("protobuf")}, true
	case mime.ProtobufJSONMediaType:
		return Media{Type: mime.ProtobufJSONMediaType, Subtype: "pbjson", Encoder: enc.Get("pbjson")}, true
	case mime.ProtobufTextMediaType:
		return Media{Type: mime.ProtobufTextMediaType, Subtype: "pbtxt", Encoder: enc.Get("pbtxt")}, true
	default:
		return Media{}, false
	}
}

func newMedia(media content.MediaType, enc *encoding.Map) Media {
	if media.Subtype == errorSubtype {
		return Media{Type: media.String(), Subtype: media.Subtype}
	}

	e := enc.Get(media.Subtype)
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: media.String(), Subtype: media.Subtype, Encoder: e}
}

func jsonMedia(enc *encoding.Map) Media {
	return Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
}
