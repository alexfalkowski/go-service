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

// NewContent returns a Content that resolves encoders from the provided map.
func NewContent(enc *encoding.Map) *Content {
	return &Content{enc: enc}
}

// Content resolves encoders from HTTP media types and provides Media helpers.
type Content struct {
	enc *encoding.Map
}

// NewFromRequest parses the request Content-Type and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
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
// If the subtype is "error" it returns a Media without an encoder. Unknown subtypes fall back to JSON.
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

// Media describes a content type and its associated encoder.
type Media struct {
	// The encoder for the media type.
	Encoder encoding.Encoder
	// The type, e.g. text.
	Type string
	// The sub type, e.g. plain.
	Subtype string
}

// IsError reports whether the media subtype represents an error payload.
func (t *Media) IsError() bool {
	return t.Subtype == errorSubtype
}
