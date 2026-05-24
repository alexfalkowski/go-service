package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

const (
	jsonKind           = "json"
	errorSubtype       = "error"
	gobSubtype         = "gob"
	messagePackSubtype = "msgpack"
)

var (
	errorType        = media.MustParse(media.Error)
	humanJSONType    = media.MustParse(media.HumanJSON)
	jsonType         = media.MustParse(media.JSON)
	messagePackType  = media.MustParse(media.MessagePack)
	protobufType     = media.MustParse(media.Protobuf)
	protobufJSONType = media.MustParse(media.ProtobufJSON)
	protobufTextType = media.MustParse(media.ProtobufText)
	textType         = media.MustParse(media.Text)
	tomlType         = media.MustParse(media.TOML)
	yamlType         = media.MustParse(media.YAML)
)

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

	value, err := media.Parse(mediaType)
	if err != nil {
		return jsonMedia(enc)
	}

	return newMedia(value, enc)
}

// Media describes an HTTP media type and its associated encoder.
//
// Type is the parsed media type. Encoder may be nil when Subtype is "error".
type Media struct {
	// Encoder is the encoder/decoder associated with the media subtype.
	Encoder encoding.Encoder

	// Type is the parsed media type.
	media.Type
}

// IsError reports whether the media subtype represents an error payload.
func (t Media) IsError() bool {
	return t.Subtype() == errorSubtype
}

// CanDecodeRequest reports whether the media type is allowed for decoding HTTP request bodies.
func (t Media) CanDecodeRequest() bool {
	subtype := t.Subtype()
	return subtype != gobSubtype && subtype != messagePackSubtype
}

// WithUTF8 returns the media type with a UTF-8 charset parameter for text media types.
func (t Media) WithUTF8() string {
	return t.Type.WithUTF8()
}

func knownMedia(mediaType string, enc *encoding.Map) (Media, bool) {
	// Exact built-in media types avoid the general parser on hot request paths.
	// Parameterized values still use the parser so their normalized Type string stays unchanged.
	switch mediaType {
	case media.Error:
		return newKnownMedia(errorType, errorSubtype, enc), true
	case media.HumanJSON:
		return newKnownMedia(humanJSONType, "hjson", enc), true
	case media.JSON:
		return jsonMedia(enc), true
	case media.MessagePack:
		return newKnownMedia(messagePackType, messagePackSubtype, enc), true
	case media.TOML:
		return newKnownMedia(tomlType, "toml", enc), true
	case media.YAML:
		return newKnownMedia(yamlType, "yaml", enc), true
	default:
		return knownProtoMedia(mediaType, enc)
	}
}

func knownProtoMedia(mediaType string, enc *encoding.Map) (Media, bool) {
	switch mediaType {
	case media.Protobuf:
		return newKnownMedia(protobufType, "protobuf", enc), true
	case media.ProtobufJSON:
		return newKnownMedia(protobufJSONType, "pbjson", enc), true
	case media.ProtobufText:
		return newKnownMedia(protobufTextType, "pbtxt", enc), true
	default:
		return Media{}, false
	}
}

func newMedia(mediaType media.Type, enc *encoding.Map) Media {
	subtype := mediaType.Subtype()
	if subtype == errorSubtype {
		return Media{Type: mediaType}
	}

	e := enc.Get(subtype)
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: mediaType, Encoder: e}
}

func newKnownMedia(mediaType media.Type, subtype string, enc *encoding.Map) Media {
	if subtype == errorSubtype {
		return Media{Type: mediaType}
	}

	e := enc.Get(subtype)
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: mediaType, Encoder: e}
}

func jsonMedia(enc *encoding.Map) Media {
	return Media{Type: jsonType, Encoder: enc.Get(jsonKind)}
}
