package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
)

// TypeKey is the HTTP header key used for Content-Type.
const TypeKey = "Content-Type"

// AcceptKey is the HTTP header key used for Accept.
const AcceptKey = "Accept"

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
// If Content-Type is not set, it falls back to the first media type in Accept.
//
// If parsing fails, it falls back to JSON.
// If the internal error media type is selected, it falls back to plain text.
func (c *Content) NewFromRequest(req *http.Request) Media {
	mediaType := req.Header.Get(TypeKey)
	if strings.IsEmpty(mediaType) {
		mediaType = firstMediaType(req.Header.Get(AcceptKey))
	}

	return c.newRequestMedia(mediaType)
}

// NewFromContentType parses the request Content-Type header and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
// If the internal error media type is selected, it falls back to plain text.
func (c *Content) NewFromContentType(req *http.Request) Media {
	return c.newRequestMedia(req.Header.Get(TypeKey))
}

// NewFromRequestBody parses the request Content-Type header and returns a matching Media for body decoding.
//
// It rejects media types that are available for internal use but intentionally unsupported for public
// request-body decoding.
func (c *Content) NewFromRequestBody(req *http.Request) (Media, error) {
	media := c.NewFromContentType(req)
	if !media.IsRequestBodySupported() {
		return media, ErrUnsupportedRequestMedia
	}

	return media, nil
}

// NewFromMedia parses mediaType and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
func (c *Content) NewFromMedia(mediaType string) Media {
	return NewMedia(mediaType, c.enc)
}

func firstMediaType(value string) string {
	mediaType, _, _ := strings.Cut(value, ",")
	return strings.TrimSpace(mediaType)
}

func (c *Content) newRequestMedia(mediaType string) Media {
	m := NewMedia(mediaType, c.enc)
	if m.IsError() {
		return newMedia(media.Text, "plain", c.enc)
	}

	return m
}
