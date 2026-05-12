package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-sync"
)

// TypeKey is the HTTP header key used for Content-Type.
const TypeKey = "Content-Type"

// NewContent constructs a Content that resolves encoders from enc and buffers responses using pool.
func NewContent(enc *encoding.Map, pool *sync.BufferPool) *Content {
	return &Content{enc: enc, pool: pool}
}

// Content resolves encoders from HTTP media types and provides helpers for content-aware request/response handling.
//
// It uses an encoding.Map registry to resolve an encoder by media subtype (e.g. "json", "hjson", "yaml", "toml").
//
// Fallback behavior:
//   - If media type parsing fails, Content falls back to JSON.
//   - If the parsed subtype is unknown (no encoder registered), Content falls back to JSON.
//
// Error subtype behavior:
//   - If the parsed subtype is "error", NewMedia returns a Media without an encoder.
//     Callers typically treat the body as a plain-text error message.
//
// Response buffering:
//   - HTTP content handlers built on this type encode successful responses into the shared buffer pool before
//     writing to the live response writer, so late encode failures do not commit partial success bodies.
type Content struct {
	enc  *encoding.Map
	pool *sync.BufferPool
}

// NewFromRequest parses the request Content-Type header and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
//
// Note: this parses the request Content-Type, not the Accept header.
func (c *Content) NewFromRequest(req *http.Request) Media {
	return NewMedia(req.Header.Get(TypeKey), c.enc)
}

// NewFromMedia parses mediaType and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
func (c *Content) NewFromMedia(mediaType string) Media {
	return NewMedia(mediaType, c.enc)
}
