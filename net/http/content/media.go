package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	content "github.com/elnormous/contenttype"
)

const jsonKind = "json"

const errorSubtype = "error"

// NewMedia builds a Media from a parsed media type and encoder map.
//
// Encoder selection:
//   - If the subtype is "error", it returns a Media without an encoder.
//   - If no encoder is registered for the subtype, it falls back to JSON.
//   - Otherwise it returns the encoder registered for the subtype.
func NewMedia(mediaType content.MediaType, enc *encoding.Map) Media {
	return newMedia(mediaType, enc)
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
	case media.JSON:
		return jsonMedia(enc), true
	case media.HJSON:
		return Media{Type: media.HJSON, Subtype: "hjson", Encoder: enc.Get("hjson")}, true
	case media.YAML:
		return Media{Type: media.YAML, Subtype: "yaml", Encoder: enc.Get("yaml")}, true
	case media.TOML:
		return Media{Type: media.TOML, Subtype: "toml", Encoder: enc.Get("toml")}, true
	case media.Protobuf:
		return Media{Type: media.Protobuf, Subtype: "protobuf", Encoder: enc.Get("protobuf")}, true
	case media.ProtobufJSON:
		return Media{Type: media.ProtobufJSON, Subtype: "pbjson", Encoder: enc.Get("pbjson")}, true
	case media.ProtobufText:
		return Media{Type: media.ProtobufText, Subtype: "pbtxt", Encoder: enc.Get("pbtxt")}, true
	default:
		return Media{}, false
	}
}

func newMedia(mediaType content.MediaType, enc *encoding.Map) Media {
	if mediaType.Subtype == errorSubtype {
		return Media{Type: mediaType.String(), Subtype: mediaType.Subtype}
	}

	e := enc.Get(mediaType.Subtype)
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: mediaType.String(), Subtype: mediaType.Subtype, Encoder: e}
}

func jsonMedia(enc *encoding.Map) Media {
	return Media{Type: media.JSON, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
}
