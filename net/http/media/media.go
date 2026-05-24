package media

import (
	"mime"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Error is the media type used for plain-text error bodies.
//
// This is intended for responses where the body is a human-readable error message.
// Note: "text/error" is not a standard IANA media type, but is used within go-service
// for consistent internal error rendering.
const Error = "text/error"

// HTML is the media type for HTML documents.
//
// This is typically used for HTML responses or debug pages.
const HTML = "text/html"

// JPEG is the media type for JPEG images.
const JPEG = "image/jpeg"

// JSON is the media type for JSON documents.
//
// This is commonly used as the Content-Type for JSON request/response bodies.
const JSON = "application/json"

// HumanJSON is the media type for HumanJSON documents.
//
// This is commonly used as the Content-Type for HumanJSON request/response bodies.
const HumanJSON = "application/hjson"

// Markdown is the media type for Markdown documents.
const Markdown = "text/markdown"

// MessagePack is the vendor media type for MessagePack payloads.
const MessagePack = "application/vnd.msgpack"

// Protobuf is the media type for protobuf binary payloads.
//
// This is commonly used when transporting protobuf wire-format bodies over HTTP.
const Protobuf = "application/protobuf"

// ProtobufJSON is the media type for protobuf JSON-encoded payloads.
//
// Note: this is a go-service specific media type string used to distinguish protobuf JSON
// from generic JSON in content negotiation.
const ProtobufJSON = "application/pbjson"

// ProtobufText is the media type for protobuf text-format payloads.
//
// Note: this is a go-service specific media type string used to distinguish protobuf text format
// in content negotiation.
const ProtobufText = "application/pbtxt"

// Text is the media type for plain text.
const Text = "text/plain"

// TOML is the media type for TOML documents.
const TOML = "application/toml"

// YAML is the media type for YAML documents.
const YAML = "application/yaml"

// ErrInvalidType is returned when a media type cannot be parsed.
var ErrInvalidType = errors.New("media: invalid type")

// Parse parses value into a Type.
func Parse(value string) (Type, error) {
	mediaType, params, err := mime.ParseMediaType(value)
	if err != nil {
		return Type{}, ErrInvalidType
	}

	_, subtype, ok := strings.Cut(mediaType, "/")
	if !ok || strings.IsEmpty(subtype) {
		return Type{}, ErrInvalidType
	}

	_, hasCharset := params["charset"]
	subtype = normalizeSubtype(subtype)

	return Type{hasCharset: hasCharset, source: value, subtype: subtype, value: mediaType}, nil
}

// MustParse parses value into a Type and panics if parsing fails.
func MustParse(value string) Type {
	mediaType, err := Parse(value)
	runtime.Must(err)

	return mediaType
}

// Type is a parsed media type.
type Type struct {
	source     string
	subtype    string
	value      string
	hasCharset bool
}

// IsZero reports whether the type has not been initialized.
func (t Type) IsZero() bool {
	return strings.IsEmpty(t.value)
}

// String returns the normalized base media type.
func (t Type) String() string {
	return t.value
}

// Subtype returns the parsed media subtype.
func (t Type) Subtype() string {
	return t.subtype
}

// WithUTF8 appends a UTF-8 charset parameter to text media types.
func (t Type) WithUTF8() string {
	if !strings.HasPrefix(t.value, "text/") {
		return t.source
	}

	if t.hasCharset {
		return t.source
	}

	return strings.Concat(t.source, "; ", "charset=utf-8")
}

// TypeByExtension returns the media type associated with ext.
//
// It wraps mime.TypeByExtension so HTTP packages can use the shared media package
// for media type lookup.
func TypeByExtension(ext string) string {
	return mime.TypeByExtension(ext)
}

func normalizeSubtype(subtype string) string {
	if subtype == "vnd.msgpack" {
		return "msgpack"
	}

	return subtype
}
