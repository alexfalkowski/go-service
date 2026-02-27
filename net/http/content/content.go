package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
	content "github.com/elnormous/contenttype"
)

const (
	jsonKind     = "json"
	errorSubtype = "error"

	// TypeKey is the HTTP header key used for Content-Type.
	TypeKey = "Content-Type"
)

// NewContent constructs a Content that resolves encoders from enc.
func NewContent(enc *encoding.Map) *Content {
	return &Content{enc: enc}
}

// Content resolves encoders from HTTP media types and provides helpers for content-aware request/response handling.
//
// It uses an encoding.Map registry to resolve an encoder by media subtype (e.g. "json", "yaml", "toml").
//
// Fallback behavior:
//   - If media type parsing fails, Content falls back to JSON.
//   - If the parsed subtype is unknown (no encoder registered), Content falls back to JSON.
//
// Error subtype behavior:
//   - If the parsed subtype is "error", NewMedia returns a Media without an encoder.
//     Callers typically treat the body as a plain-text error message.
type Content struct {
	enc *encoding.Map
}

// NewFromRequest parses the request Content-Type header and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
//
// Note: this parses the request Content-Type, not the Accept header.
func (c *Content) NewFromRequest(req *http.Request) *Media {
	t, err := content.GetMediaType(req)
	if err != nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: c.enc.Get(jsonKind)}
	}

	return NewMedia(t, c.enc)
}

// NewFromMedia parses mediaType and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
func (c *Content) NewFromMedia(mediaType string) *Media {
	t, err := content.ParseMediaType(mediaType)
	if err != nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: c.enc.Get(jsonKind)}
	}

	return NewMedia(t, c.enc)
}

// NewMedia builds a Media from a parsed media type and encoder map.
//
// Encoder selection:
//   - If the subtype is "error", it returns a Media without an encoder.
//   - If no encoder is registered for the subtype, it falls back to JSON.
//   - Otherwise it returns the encoder registered for the subtype.
func NewMedia(media content.MediaType, enc *encoding.Map) *Media {
	if media.Subtype == errorSubtype {
		return &Media{Type: media.String(), Subtype: media.Subtype}
	}

	e := enc.Get(media.Subtype)
	if e == nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	return &Media{Type: media.String(), Subtype: media.Subtype, Encoder: e}
}

// Media describes an HTTP media type and its associated encoder.
//
// Type is the full media type string (including parameters if present). Subtype is the parsed media subtype used
// for encoder lookup. Encoder may be nil when Subtype is "error".
type Media struct {
	// Encoder is the encoder/decoder associated with the media subtype.
	Encoder encoding.Encoder

	// Type is the full media type string (for example "application/json" or "text/plain; charset=utf-8").
	Type string

	// Subtype is the parsed subtype (for example "json" or "plain").
	Subtype string
}

// IsError reports whether the media subtype represents an error payload.
func (t *Media) IsError() bool {
	return t.Subtype == errorSubtype
}
