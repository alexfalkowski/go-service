package unary

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http/content/policy"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

const (
	jsonKind     = "json"
	errorSubtype = "error"
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

// unaryKinds maps a media subtype alias to the [encoding.Map] kind that actually encodes/decodes it,
// mirroring how streamKinds in stream_media.go maps a subtype to a [stream.Map] kind. encoding.Map
// registers each encoder under exactly one canonical kind (see [encoding.NewMap]), so every other
// spelling HTTP clients may send has to be translated here.
//
// Unlike streamKinds, this is deliberately not an allowlist: a subtype absent from this map still
// resolves by using the subtype as the kind directly (see unaryKind), because most media subtypes
// already match their encoding.Map kind by name (see
// TestNewFromRequestBodyDecodesFallthroughReachableMediaTypes). The entries below exist only for the
// subtypes that do not.
var unaryKinds = map[string]string{
	"pb":           "protobuf",
	"pbbin":        "protobuf",
	"proto":        "protobuf",
	"protobin":     "protobuf",
	"pbtxt":        "prototext",
	"prototxt":     "prototext",
	"pbjson":       "protojson",
	"octet-stream": "bytes",
	"plain":        "bytes",
	"yml":          "yaml",
}

// unaryKind resolves subtype to its registered encoding.Map kind, defaulting to subtype itself when no
// explicit entry exists in unaryKinds.
func unaryKind(subtype string) string {
	if kind, ok := unaryKinds[subtype]; ok {
		return kind
	}

	return subtype
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
//
// A nil Encoder means the media type resolved to no codec at all, which is never decodable. Otherwise
// the codec must satisfy the decoder-bounds rule; see undecodableKinds.
//
// This answers for the media type this Media actually holds. A Media produced by the outbound JSON
// fallback therefore reports JSON's answer, because JSON is genuinely what it holds; only
// Content.NewFromRequestBody resolves a request Content-Type strictly.
func (t Media) CanDecodeRequest() bool {
	return t.Encoder != nil && policy.CanDecode(t.Subtype())
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
		return newMedia(errorType, enc), true
	case media.HumanJSON:
		return newMedia(humanJSONType, enc), true
	case media.JSON:
		return jsonMedia(enc), true
	case media.MessagePack:
		return newMedia(messagePackType, enc), true
	case media.TOML:
		return newMedia(tomlType, enc), true
	case media.YAML:
		return newMedia(yamlType, enc), true
	default:
		return knownProtoMedia(mediaType, enc)
	}
}

func knownProtoMedia(mediaType string, enc *encoding.Map) (Media, bool) {
	switch mediaType {
	case media.Protobuf:
		return newMedia(protobufType, enc), true
	case media.ProtobufJSON:
		return newMedia(protobufJSONType, enc), true
	case media.ProtobufText:
		return newMedia(protobufTextType, enc), true
	default:
		return Media{}, false
	}
}

func newMedia(mediaType media.Type, enc *encoding.Map) Media {
	subtype := mediaType.Subtype()
	if subtype == errorSubtype {
		return Media{Type: mediaType}
	}

	e := enc.Get(unaryKind(subtype))
	if e == nil {
		return jsonMedia(enc)
	}

	return Media{Type: mediaType, Encoder: e}
}

func jsonMedia(enc *encoding.Map) Media {
	return Media{Type: jsonType, Encoder: enc.Get(jsonKind)}
}
