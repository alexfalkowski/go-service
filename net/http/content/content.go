package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonMediaType = "application/json"
	jsonKind      = "json"

	// TypeKey for HTTP headers.
	TypeKey = "Content-Type"
)

// Decode to v, unless an error occurs.
func Decode(ctx context.Context, v any) error {
	e := hc.Encoder(ctx)
	req := hc.Request(ctx)

	return e.Decode(req.Body, v)
}

// Content creates types from media types.
type Content struct {
	enc *encoding.Map
}

// NewContent with an encoding.
func NewContent(enc *encoding.Map) *Content {
	return &Content{enc: enc}
}

// NewFromRequest for content.
func (c *Content) NewFromRequest(req *http.Request) *Type {
	t, err := ct.GetMediaType(req)

	return newType(t, err, c.enc)
}

// NewFromMedia for content.
func (c *Content) NewFromMedia(mediaType string) *Type {
	t, err := ct.ParseMediaType(mediaType)

	return newType(t, err, c.enc)
}

// Type for content.
type Type struct {
	Encoder encoding.Encoder
	Media   string
	Kind    string
}

// IsText for type.
func (t *Type) IsText() bool {
	return t.Kind == "plain"
}

func newType(t ct.MediaType, err error, enc *encoding.Map) *Type {
	if err != nil {
		return &Type{Media: jsonMediaType, Kind: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	if t.Subtype == "plain" {
		return &Type{Media: t.String(), Kind: t.Subtype}
	}

	e := enc.Get(t.Subtype)
	if e == nil {
		return &Type{Media: jsonMediaType, Kind: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	return &Type{Media: t.String(), Kind: t.Subtype, Encoder: e}
}
