package stream

import (
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content/policy"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-sync"
)

// NewContent constructs Content that resolves streaming codecs from enc and buffers initial responses using pool.
func NewContent(enc *stream.Map, pool *sync.BufferPool) *Content {
	return &Content{enc: enc, pool: pool}
}

// Content resolves streaming codecs from HTTP media types and buffers initial responses before they commit.
//
// It uses an [stream.Map] registry to resolve streaming encoders and decoders. Unlike unary content,
// streaming request and response media negotiation is strict: unsupported declared media is rejected rather
// than falling back to another representation.
type Content struct {
	enc  *stream.Map
	pool *sync.BufferPool
}

// NewFromMedia builds a Media from mediaType and Content's streaming registry.
//
// NewFromMedia never falls back to JSON: an unparsable media type, a media type with
// no entry in the streaming media mapping, or an unregistered kind all return
// [ErrUnsupportedMedia]. A nil Content also returns [ErrUnsupportedMedia]. Streaming callers must fail
// explicitly rather than silently degrading to a different wire format.
func (c *Content) NewFromMedia(mediaType string) (Media, error) {
	if c == nil {
		return Media{}, ErrUnsupportedMedia
	}

	if resolved, ok := knownMedia(mediaType, c.enc); ok {
		return usableMedia(resolved)
	}

	value, err := media.Parse(mediaType)
	if err != nil {
		return Media{}, ErrUnsupportedMedia
	}

	kind, ok := streamKinds[value.String()]
	if !ok {
		return Media{}, ErrUnsupportedMedia
	}

	codec := c.enc.Get(kind)

	return usableMedia(Media{NewEncoder: codec.Encoder, NewDecoder: codec.Decoder, Type: value})
}

// NewFromAccept resolves the request Accept header to a streaming encoder, falling back to
// Content-Type when Accept is absent, and to [media.NDJSON] when neither is set. A registered
// codec without an encoder is unsupported. A nil Content returns [ErrUnsupportedMedia].
//
// An Accept list is satisfiable for the server's producible streaming media type if it contains an
// exact match without non-quality parameters, or a matching wildcard ("*/*", or "type/*" where type
// matches the producible type's own major type) without non-quality parameters, anywhere in the list.
// The canonical NDJSON representation has no media parameters, so a range that requires one cannot match it.
// Per RFC 9110 §12.5.1, the most specific reference present controls
// regardless of list order: an exact subtype match takes precedence over any wildcard, and a "type/*"
// wildcard takes precedence over the bare "*/*" wildcard. Only that controlling reference's explicit
// q=0 exclusion (see [github.com/alexfalkowski/go-service/v2/net/http/accept.IsZeroQuality]) decides
// satisfiability — a q=0 on a less specific reference elsewhere in the list has no effect once a more
// specific reference is present, and conversely a q=0 on the controlling reference cannot be overridden
// by a less specific, non-excluded reference. A satisfiable list resolves to its exact match if the
// controlling reference was one, otherwise to media.NDJSON.
//
// Unlike unary response-media negotiation, a non-empty Accept list that is not satisfiable this way returns
// [ErrUnsupportedMedia] instead of falling back to JSON or to Content-Type: a client that named
// only concrete, unproducible media types must get an explicit rejection, not a silently different wire
// format.
func (c *Content) NewFromAccept(req *http.Request) (Media, error) {
	if c == nil {
		return Media{}, ErrUnsupportedMedia
	}

	accept := req.Header.Values(http.AcceptKey)
	if !strings.IsEmpty(strings.Join(strings.Empty, accept...)) {
		return c.newFromAcceptHeader(strings.Join(", ", accept...))
	}

	mediaType := req.Header.Get(http.ContentTypeKey)
	if strings.IsEmpty(mediaType) {
		mediaType = media.NDJSON
	}

	return c.newEncoderMedia(mediaType)
}

// NewFromContentType parses the request Content-Type header and resolves a streaming decoder.
//
// Unlike unary request-media negotiation, an unregistered or unparsable media type returns
// [ErrUnsupportedMedia] instead of falling back to JSON.
//
// This is the streaming counterpart of unary request-media negotiation: after NewFromMedia resolves the
// media type, it re-checks the resolved kind against the same request-decode policy that guards
// unary request bodies (see the decoder-bounds rule in the package documentation), so a streaming media
// type that maps to a banned kind is rejected here even though NewFromMedia itself has no opinion on
// that policy. NewFromAccept and NewFromMedia are not policed this way: only this method decodes an
// untrusted request body (see [Media.NewDecoder]'s other caller, the client's response-stream path,
// which must keep decoding every registered kind).
func (c *Content) NewFromContentType(req *http.Request) (Media, error) {
	m, err := c.NewFromMedia(req.Header.Get(http.ContentTypeKey))
	if err != nil {
		return Media{}, err
	}

	// Re-resolve the kind rather than carrying it on Media, keeping the struct free of a field
	// that exists only for this one caller. A missing entry, a nil decoder constructor, or a banned kind
	// are all treated as unsupported so this fails closed: neither the missing-entry nor the banned-kind
	// arm can currently fire, because streamKinds' only mapping ("application/x-ndjson" to "json") already passed the
	// lookup inside NewFromMedia and "json" is not a banned kind. Both stay as guards for a future
	// streamKinds entry that maps to a banned kind, or a NewFromMedia change that admits an unknown media
	// type.
	kind, ok := streamKinds[m.String()]
	if !ok || m.NewDecoder == nil || !policy.CanDecode(kind) {
		return Media{}, ErrUnsupportedMedia
	}

	return m, nil
}

func (c *Content) newFromAcceptHeader(header string) (Media, error) {
	if resolved, ok := knownMedia(header, c.enc); ok {
		return encoderMedia(resolved, nil)
	}

	mediaType, ok := matchStreamAccept(header)
	if !ok {
		return Media{}, ErrUnsupportedMedia
	}

	return c.newEncoderMedia(mediaType)
}

func (c *Content) newEncoderMedia(mediaType string) (Media, error) {
	return encoderMedia(c.NewFromMedia(mediaType))
}
