package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

const jsonKind = "json"

const errorSubtype = "error"

// NewMedia builds a Media from a media type string and encoder map.
//
// Encoder selection:
//   - If the subtype is "error", it returns a Media without an encoder.
//   - If no encoder is registered for the subtype, it falls back to JSON.
//   - Otherwise it returns the encoder registered for the subtype.
func NewMedia(mediaType string, enc *encoding.Map) Media {
	if media, ok := knownMedia(mediaType, enc); ok {
		return media
	}

	value, subtype, err := media.Parse(mediaType)
	if err != nil {
		return jsonMedia(enc)
	}

	return newMedia(value, subtype, enc)
}

// Media describes an HTTP media type and its associated encoder.
//
// Type is the base media type string. Subtype is the parsed media subtype used for encoder lookup.
// Encoder may be nil when Subtype is "error".
type Media struct {
	// Encoder is the encoder/decoder associated with the media subtype.
	Encoder encoding.Encoder

	// Type is the base media type string (for example "application/json", "application/hjson", or "text/plain").
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

func newMedia(value, subtype string, enc *encoding.Map) Media {
	if subtype == errorSubtype {
		return Media{Type: value, Subtype: subtype}
	}

	e := enc.Get(subtype)
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: value, Subtype: subtype, Encoder: e}
}

func jsonMedia(enc *encoding.Map) Media {
	return Media{Type: media.JSON, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
}
