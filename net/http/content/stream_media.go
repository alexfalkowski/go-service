package content

import (
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// streamKinds maps a media subtype to the [stream.Map] kind that encodes/decodes it.
//
// This mapping is explicit rather than reusing the subtype as the kind directly: "application/x-ndjson"
// has subtype "x-ndjson", but it is framed newline-delimited JSON, so it reuses the "json" streaming
// kind rather than needing its own codec.
var streamKinds = map[string]string{
	"x-ndjson": jsonKind,
}

// NewStreamMedia builds a StreamMedia from a media type string and streaming registry.
//
// Unlike [NewMedia], NewStreamMedia never falls back to JSON: an unparsable media type, a subtype with
// no entry in the streaming media mapping, or a kind unregistered in sm all return
// [ErrUnsupportedStreamMedia]. Streaming callers must fail explicitly rather than silently degrading to
// a different wire format.
func NewStreamMedia(mediaType string, sm *stream.Map) (StreamMedia, error) {
	value, err := media.Parse(mediaType)
	if err != nil {
		return StreamMedia{}, ErrUnsupportedStreamMedia
	}

	kind, ok := streamKinds[value.Subtype()]
	if !ok {
		return StreamMedia{}, ErrUnsupportedStreamMedia
	}

	newEncoder := sm.GetEncoder(kind)
	newDecoder := sm.GetDecoder(kind)
	if newEncoder == nil && newDecoder == nil {
		return StreamMedia{}, ErrUnsupportedStreamMedia
	}

	return StreamMedia{NewEncoder: newEncoder, NewDecoder: newDecoder, Type: value}, nil
}

// StreamMedia describes a resolved streaming media type and its encoder/decoder constructors.
//
// NewEncoder or NewDecoder may be nil when the registered kind supports only one direction; callers
// must check the one they intend to use.
type StreamMedia struct {
	// NewEncoder constructs a [stream.Encoder] bound to a response writer, or nil if the kind has no
	// registered encoder constructor.
	NewEncoder stream.EncoderFunc

	// NewDecoder constructs a [stream.Decoder] bound to a request body, or nil if the kind has no
	// registered decoder constructor.
	NewDecoder stream.DecoderFunc

	// Type is the parsed media type.
	media.Type
}

// NewStreamFromAccept parses the first request Accept media type and resolves a streaming encoder.
//
// If Accept is not set, it falls back to Content-Type. Unlike [Content.NewFromAccept], an unregistered
// or unparsable media type returns [ErrUnsupportedStreamMedia] instead of falling back to JSON.
func (c *Content) NewStreamFromAccept(req *http.Request) (StreamMedia, error) {
	mediaType := firstListItem(req.Header.Get(AcceptKey))
	if strings.IsEmpty(mediaType) {
		mediaType = req.Header.Get(TypeKey)
	}

	return NewStreamMedia(mediaType, c.sm)
}

// NewStreamFromContentType parses the request Content-Type header and resolves a streaming decoder.
//
// Unlike [Content.NewFromContentType], an unregistered or unparsable media type returns
// [ErrUnsupportedStreamMedia] instead of falling back to JSON.
//
// This is the streaming counterpart of [Content.NewFromRequestBody]: after NewStreamMedia resolves the
// media type, it re-checks the resolved kind against the same request-decode policy that guards
// unary request bodies (see the decoder-bounds rule in the package documentation), so a streaming media
// type that maps to a banned kind is rejected here even though NewStreamMedia itself has no opinion on
// that policy. NewStreamFromAccept and [Content.NewStreamFromMedia] are not policed this way: only this
// constructor decodes an untrusted request body (see [StreamMedia.NewDecoder]'s other caller, the
// client's response-stream path, which must keep decoding every registered kind).
func (c *Content) NewStreamFromContentType(req *http.Request) (StreamMedia, error) {
	m, err := NewStreamMedia(req.Header.Get(TypeKey), c.sm)
	if err != nil {
		return StreamMedia{}, err
	}

	// Re-resolve the kind rather than carrying it on StreamMedia, keeping the struct free of a field
	// that exists only for this one caller. A missing entry, a nil decoder constructor, or a banned kind
	// are all treated as unsupported so this fails closed: neither the missing-entry nor the banned-kind
	// arm can currently fire, because streamKinds' only mapping ("x-ndjson" to "json") already passed the
	// lookup inside NewStreamMedia and "json" is not a banned kind. Both stay as guards for a future
	// streamKinds entry that maps to a banned kind, or a NewStreamMedia change that admits an unknown
	// subtype.
	kind, ok := streamKinds[m.Subtype()]
	if !ok || m.NewDecoder == nil || !canDecodeRequest(kind) {
		return StreamMedia{}, ErrUnsupportedStreamMedia
	}

	return m, nil
}
