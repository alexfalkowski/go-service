package media

import (
	"github.com/alexfalkowski/go-service/v2/errors"
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

// HJSON is the media type for HJSON documents.
//
// This is commonly used as the Content-Type for HJSON request/response bodies.
const HJSON = "application/hjson"

// Markdown is the media type for Markdown documents.
const Markdown = "text/markdown"

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

// Parse parses value into a base media type and subtype.
//
// Parameters are ignored because content negotiation only uses the base media type.
func Parse(value string) (string, string, error) {
	mediaType, _, _ := strings.Cut(value, ";")
	mediaType = strings.TrimSpace(mediaType)

	_, subtype, ok := strings.Cut(mediaType, "/")
	if !ok || strings.IsEmpty(subtype) {
		return strings.Empty, strings.Empty, ErrInvalidType
	}

	return mediaType, subtype, nil
}

// WithUTF8 appends a UTF-8 charset parameter to text media types.
//
// Non-text media types and media types that already contain a charset parameter are returned unchanged.
func WithUTF8(mediaType string) string {
	value, _, err := Parse(mediaType)
	if err != nil || !strings.HasPrefix(value, "text/") || strings.Contains(mediaType, "charset=") {
		return mediaType
	}

	return strings.Concat(mediaType, "; ", "charset=utf-8")
}
