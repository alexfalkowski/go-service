package content

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonMediaType = "application/json"
	jsonKind      = "json"
	plainSubtype  = "plain"

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

	return newType(t, err, c.enc)
}

// NewFromMedia for content.
func (c *Content) NewFromMedia(mediaType string) *Media {
	t, err := ct.ParseMediaType(mediaType)

	return newType(t, err, c.enc)
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

// IsText for type.
func (t *Media) IsText() bool {
	return t.Subtype == plainSubtype
}

func newType(media ct.MediaType, err error, enc *encoding.Map) *Media {
	if err != nil {
		return &Media{Type: jsonMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	if media.Subtype == plainSubtype {
		return &Media{Type: media.String(), Subtype: media.Subtype}
	}

	e := enc.Get(media.Subtype)

	return &Media{Type: media.String(), Subtype: media.Subtype, Encoder: e}
}
