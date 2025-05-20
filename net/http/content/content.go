package content

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/mime"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonKind     = "json"
	errorSubtype = "error"

	// TypeKey for HTTP headers.
	TypeKey = "Content-Type"
)

// NewContent with an encoding.
func NewContent(enc *encoding.Map) *Content {
	return &Content{enc: enc}
}

// Content creates types from media types.
type Content struct {
	enc *encoding.Map
}

// NewFromRequest for content.
func (c *Content) NewFromRequest(req *http.Request) *Media {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: c.enc.Get(jsonKind)}
	}

	return NewMedia(t, c.enc)
}

// NewFromMedia for content.
func (c *Content) NewFromMedia(mediaType string) *Media {
	t, err := ct.ParseMediaType(mediaType)
	if err != nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: c.enc.Get(jsonKind)}
	}

	return NewMedia(t, c.enc)
}

// NewMedia for content.
func NewMedia(media ct.MediaType, enc *encoding.Map) *Media {
	if media.Subtype == errorSubtype {
		return &Media{Type: media.String(), Subtype: media.Subtype}
	}

	e := enc.Get(media.Subtype)
	if e == nil {
		return &Media{Type: mime.JSONMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	return &Media{Type: media.String(), Subtype: media.Subtype, Encoder: e}
}

// Media for content.
// https://en.wikipedia.org/wiki/Media_type
type Media struct {
	// The encoder for the media type.
	Encoder encoding.Encoder

	// The type, e.g. text.
	Type string

	// The sub type, e.g. plain.
	Subtype string
}

// IsError for type.
func (t *Media) IsError() bool {
	return t.Subtype == errorSubtype
}
