package content

import (
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/sync"
	ct "github.com/elnormous/contenttype"
)

const (
	jsonMediaType = "application/json"
	jsonKind      = "json"

	// TypeKey for HTTP headers.
	TypeKey = "Content-Type"
)

var pool *sync.BufferPool

// Register for content.
func Register(p *sync.BufferPool) {
	pool = p
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
func (c *Content) NewFromRequest(req *http.Request) *Media {
	t, err := ct.GetMediaType(req)

	return newType(t, err, c.enc)
}

// NewFromMedia for content.
func (c *Content) NewFromMedia(mediaType string) *Media {
	t, err := ct.ParseMediaType(mediaType)

	return newType(t, err, c.enc)
}

func (c *Content) handler(prefix string, handler handler) http.HandlerFunc {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix(prefix, runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := c.NewFromRequest(req)

		ctx = hc.WithEncoder(ctx, ct.Encoder)
		res.Header().Add(TypeKey, ct.Type)

		data, err := handler(ctx)
		runtime.Must(err)

		err = ct.Encoder.Encode(res, data)
		runtime.Must(err)
	}

	return h
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
	return t.Subtype == "plain"
}

func newType(t ct.MediaType, err error, enc *encoding.Map) *Media {
	if err != nil {
		return &Media{Type: jsonMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	if t.Subtype == "plain" {
		return &Media{Type: t.String(), Subtype: t.Subtype}
	}

	e := enc.Get(t.Subtype)
	if e == nil {
		return &Media{Type: jsonMediaType, Subtype: jsonKind, Encoder: enc.Get(jsonKind)}
	}

	return &Media{Type: t.String(), Subtype: t.Subtype, Encoder: e}
}
