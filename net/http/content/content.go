package content

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonMediaType = "application/json"
	jsonKind      = "json"

	// TypeKey for HTTP headers.
	TypeKey = "Content-Type"
)

// NewFromRequest for content.
func NewFromRequest(req *http.Request, enc *encoding.Map) *Type {
	t, err := ct.GetMediaType(req)

	return newType(t, err, enc)
}

// NewFromMedia for content.
func NewFromMedia(mediaType string, enc *encoding.Map) *Type {
	t, err := ct.ParseMediaType(mediaType)

	return newType(t, err, enc)
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
