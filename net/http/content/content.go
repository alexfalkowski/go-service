package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
)

// TypeKey is the HTTP header key used for Content-Type.
const TypeKey = "Content-Type"

// AcceptKey is the HTTP header key used for Accept.
const AcceptKey = "Accept"

// NewContent constructs a Content that resolves single-value encoders from enc, streaming encoders/decoders
// from sm, and buffers responses using pool.
func NewContent(enc *encoding.Map, sm *stream.Map, pool *sync.BufferPool) *Content {
	return &Content{enc: enc, sm: sm, pool: pool}
}

// Content resolves encoders from HTTP media types and provides helpers for content-aware request/response handling.
//
// It uses an [encoding.Map] registry to resolve a single-value encoder by media subtype (e.g. "json", "hjson",
// "yaml", "toml"), and a [stream.Map] registry to resolve a streaming encoder/decoder by media subtype for
// streaming routes (e.g. "x-ndjson"). The two registries are intentionally separate: an unregistered streaming
// kind must fail explicitly rather than silently falling back to single-value JSON (see [NewStreamMedia]).
//
// Fallback behavior:
//   - If media type parsing fails, Content falls back to JSON.
//   - If the parsed subtype is unknown (no encoder registered), Content falls back to JSON.
//   - This fallback is deliberately outbound-only: it drives Accept-based response encoder selection
//     and Content.NewFromMedia. [Content.NewFromRequestBody] resolves an inbound request Content-Type
//     strictly instead — an absent Content-Type still defaults to JSON, but an unparseable or
//     unregistered one is rejected rather than silently reinterpreted as a format the caller did not
//     send.
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
	sm   *stream.Map
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
		mediaType = firstListItem(req.Header.Get(AcceptKey))
	}

	return c.newRequestMedia(mediaType)
}

// NewFromAccept parses the first request Accept media type and returns a matching Media.
//
// If Accept is not set, it falls back to Content-Type.
//
// If parsing fails, it falls back to JSON.
// If the internal error media type is selected, it falls back to plain text.
func (c *Content) NewFromAccept(req *http.Request) Media {
	mediaType := firstListItem(req.Header.Get(AcceptKey))
	if strings.IsEmpty(mediaType) {
		mediaType = req.Header.Get(TypeKey)
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
// Unlike [Content.NewFromContentType], there is no JSON fallback for an unparseable or unregistered
// media type: answering a Content-Type the caller declared with a different codec would silently decode
// the body as a format it did not send, so that case returns [ErrUnsupportedRequestMedia] instead. An
// absent Content-Type still defaults to JSON, since the caller asserted nothing.
//
// It also rejects media types that are available for internal use but intentionally unsupported for
// public request-body decoding; see the decoder-bounds rule in the package documentation.
func (c *Content) NewFromRequestBody(req *http.Request) (Media, error) {
	m, err := c.requestMedia(req.Header.Get(TypeKey))
	if err != nil {
		return Media{}, err
	}

	// text/error is response media with no codec of its own; decode its body as text, matching the
	// existing behaviour asserted by TestNewRequestHandlerTreatsInternalErrorContentTypeAsText. Resolve
	// the text codec directly rather than through newKnownMedia, whose nil-codec path would fall back to
	// JSON; a nil encoder here is rejected by CanDecodeRequest below.
	if m.IsError() {
		m = Media{Type: textType, Encoder: c.enc.Get("plain")}
	}

	if !m.CanDecodeRequest() {
		return Media{}, ErrUnsupportedRequestMedia
	}

	return m, nil
}

// requestMedia resolves a request Content-Type strictly: unlike NewMedia there is no JSON fallback for
// an unparseable or unregistered media type, because answering a declared Content-Type with a different
// codec silently decodes the body as a format the caller did not send.
//
// A returned Media may still carry a nil Encoder; the caller rejects that.
func (c *Content) requestMedia(mediaType string) (Media, error) {
	// An absent Content-Type asserts nothing, so the documented JSON default stands.
	if strings.IsEmpty(mediaType) {
		return jsonMedia(c.enc), nil
	}

	// knownMedia falls back to JSON when a known media type's codec was nil-registered. Compare the
	// resolved type against the requested one to reject that rather than decode as JSON.
	if m, ok := knownMedia(mediaType, c.enc); ok {
		if m.String() != mediaType {
			return Media{}, ErrUnsupportedRequestMedia
		}

		return m, nil
	}

	value, err := media.Parse(mediaType)
	if err != nil {
		return Media{}, ErrUnsupportedRequestMedia
	}

	// The error subtype has no codec, so recognise it before the lookup, matching newMedia. Otherwise a
	// parameterized text/error would be rejected while the exact form is coerced to text.
	if value.Subtype() == errorSubtype {
		return Media{Type: value}, nil
	}

	encoder := c.enc.Get(value.Subtype())
	if encoder == nil {
		return Media{}, ErrUnsupportedRequestMedia
	}

	return Media{Type: value, Encoder: encoder}, nil
}

// NewFromMedia parses mediaType and returns a matching Media.
//
// If parsing fails, it falls back to JSON.
func (c *Content) NewFromMedia(mediaType string) Media {
	return NewMedia(mediaType, c.enc)
}

// NewStreamFromMedia parses mediaType and resolves a streaming encoder/decoder pair, the streaming
// counterpart of NewFromMedia.
//
// Unlike NewFromMedia, there is no JSON fallback: an unparseable or unregistered streaming media type
// returns [ErrUnsupportedStreamMedia] (see [NewStreamMedia]), including when c has no streaming
// registry configured — that is indistinguishable from "this kind isn't registered" to the caller, and
// is handled the same way rather than through a separate error.
//
// This is the entry point a caller without an *[http.Request] to negotiate from (for example
// [github.com/alexfalkowski/go-service/v2/net/http/client.Client]) uses to resolve a streaming media
// type directly, mirroring how [Content.NewFromMedia] already does this for single-value media.
func (c *Content) NewStreamFromMedia(mediaType string) (StreamMedia, error) {
	return NewStreamMedia(mediaType, c.sm)
}

func (c *Content) newRequestMedia(mediaType string) Media {
	m := NewMedia(mediaType, c.enc)
	if m.IsError() {
		return newKnownMedia(textType, "plain", c.enc)
	}

	return m
}
