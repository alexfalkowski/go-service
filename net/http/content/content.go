package content

import (
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonMediaType = "application/json"
	jsonKind      = "json"
)

var (
	// TypeKey for HTTP headers.
	TypeKey = "Content-Type"

	// ErrInvalidEncoder for content.
	ErrInvalidEncoder = errors.New("invalid encoder")
)

// AddJSONHeader for content.
func AddJSONHeader(header http.Header) {
	header.Set(TypeKey, jsonMediaType)
}

// NewFromRequest for content.
func NewFromRequest(req *http.Request) *Type {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return &Type{Media: jsonMediaType, Kind: jsonKind}
	}

	return &Type{Media: t.String(), Kind: t.Subtype}
}

// NewFromMedia for content.
func NewFromMedia(mediaType string) *Type {
	t, err := ct.ParseMediaType(mediaType)
	if err != nil {
		return &Type{Media: jsonMediaType, Kind: jsonKind}
	}

	return &Type{Media: t.String(), Kind: t.Subtype}
}

// Type for content.
type Type struct {
	Media string
	Kind  string
}

// Marshaller for type.
func (t *Type) Encoder(enc *encoding.Map) (encoding.Encoder, error) {
	m := enc.Get(t.Kind)
	if m == nil {
		return nil, ErrInvalidEncoder
	}

	return m, nil
}

// IsText for type.
func (t *Type) IsText() bool {
	return t.Kind == "plain"
}
