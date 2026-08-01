package stream

import (
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/accept"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/content/policy"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ndjsonType is the parsed form of [media.NDJSON], the only streaming media type this package currently
// produces. matchStreamAccept uses its major type to decide whether an Accept wildcard is satisfiable.
var ndjsonType = media.MustParse(media.NDJSON)

// streamKinds maps a normalized media type to the [stream.Map] kind that encodes/decodes it.
//
// This mapping is explicit rather than reusing the media type as the kind directly:
// "application/x-ndjson" is framed newline-delimited JSON, so it reuses the "json" streaming kind rather
// than needing its own codec.
var streamKinds = map[string]string{
	media.NDJSON: "json",
}

// NewMedia builds a Media from a media type string and streaming registry.
//
// NewMedia never falls back to JSON: an unparsable media type, a media type with
// no entry in the streaming media mapping, or a kind unregistered in sm all return
// [ErrUnsupportedMedia]. Streaming callers must fail explicitly rather than silently degrading to
// a different wire format.
func NewMedia(mediaType string, sm *stream.Map) (Media, error) {
	value, err := media.Parse(mediaType)
	if err != nil {
		return Media{}, ErrUnsupportedMedia
	}

	kind, ok := streamKinds[value.String()]
	if !ok {
		return Media{}, ErrUnsupportedMedia
	}

	codec := sm.Get(kind)
	if codec.Encoder == nil && codec.Decoder == nil {
		return Media{}, ErrUnsupportedMedia
	}

	return Media{NewEncoder: codec.Encoder, NewDecoder: codec.Decoder, Type: value}, nil
}

// Media describes a resolved streaming media type and its encoder/decoder constructors.
//
// NewEncoder or NewDecoder may be nil when the registered kind supports only one direction; callers
// must check the one they intend to use.
type Media struct {
	// NewEncoder constructs a [stream.Encoder] bound to a response writer, or nil if the kind has no
	// registered encoder constructor.
	NewEncoder stream.EncoderFunc

	// NewDecoder constructs a [stream.Decoder] bound to a request body, or nil if the kind has no
	// registered decoder constructor.
	NewDecoder stream.DecoderFunc

	// Type is the parsed media type.
	media.Type
}

// NewMediaFromAccept resolves the request Accept header to a streaming encoder, falling back to
// Content-Type when Accept is absent, and to [media.NDJSON] when neither is set. A registered
// codec without an encoder is unsupported.
//
// An Accept list is satisfiable for the server's producible streaming media type if it contains an
// exact match, or a matching wildcard ("*/*", or "type/*" where type matches the producible type's own
// major type), anywhere in the list. Per RFC 9110 §12.5.1, the most specific reference present controls
// regardless of list order: an exact subtype match takes precedence over any wildcard, and a "type/*"
// wildcard takes precedence over the bare "*/*" wildcard. Only that controlling reference's explicit
// q=0 exclusion (see [github.com/alexfalkowski/go-service/v2/net/http/accept.IsZeroQuality]) decides
// satisfiability — a q=0 on a less specific reference elsewhere in the list has no effect once a more
// specific reference is present, and conversely a q=0 on the controlling reference cannot be overridden
// by a less specific, non-excluded reference. A satisfiable list resolves to its exact match if the
// controlling reference was one, otherwise to media.NDJSON.
//
// Unlike [content.Content.NewMediaFromAccept], a non-empty Accept list that is not satisfiable this way returns
// [ErrUnsupportedMedia] instead of falling back to JSON or to Content-Type: a client that named
// only concrete, unproducible media types must get an explicit rejection, not a silently different wire
// format.
func NewMediaFromAccept(req *http.Request, sm *stream.Map) (Media, error) {
	var mediaType string

	header := req.Header.Get(content.AcceptKey)
	if strings.IsEmpty(header) {
		m := req.Header.Get(content.TypeKey)
		if strings.IsEmpty(m) {
			m = media.NDJSON
		}

		mediaType = m
	} else {
		m, ok := matchStreamAccept(header)
		if !ok {
			return Media{}, ErrUnsupportedMedia
		}

		mediaType = m
	}

	resolved, err := NewMedia(mediaType, sm)
	if err != nil || resolved.NewEncoder == nil {
		return Media{}, ErrUnsupportedMedia
	}

	return resolved, nil
}

// matchStreamAccept reports whether header, an Accept header value, is satisfiable for a registered
// streamKinds entry, and if so, which media type to resolve.
//
// It finds the most specific reference to the producible type present anywhere in the list — an exact
// subtype match, else a "type/*" wildcard [accept.IsWildcard] reports as satisfied by [ndjsonType]
// ("text/*" does not satisfy an "application/x-ndjson" route the way "*/*" or "application/*" does),
// else the bare "*/*" wildcard — and returns satisfiable only if that one reference is not
// [accept.IsZeroQuality]; a less specific reference's own quality, zero or not, is irrelevant once a
// more specific reference is found. This matches RFC 9110 §12.5.1's "most specific reference" rule
// regardless of list order. When satisfiable, it returns the exact match's media type if the
// controlling reference was one, otherwise [media.NDJSON]. An unparsable item is skipped rather than
// rejecting the whole list, the same way an unparsable single Accept value already falls through to
// [NewMedia]'s own rejection when nothing else in the list matches.
func matchStreamAccept(header string) (string, bool) {
	var exact, major, bare streamAcceptMatch

	for _, item := range accept.Items(header) {
		value, err := media.Parse(item)
		if err != nil {
			continue
		}

		zero := accept.IsZeroQuality(item)

		if _, ok := streamKinds[value.String()]; ok {
			exact.consider(zero, value.String())
			continue
		}

		if !accept.IsWildcard(value, ndjsonType) {
			continue
		}

		if value.Major() == "*" {
			bare.consider(zero, media.NDJSON)
		} else {
			major.consider(zero, media.NDJSON)
		}
	}

	switch {
	case exact.found:
		return exact.value, !exact.zero
	case major.found:
		return major.value, !major.zero
	case bare.found:
		return bare.value, !bare.zero
	default:
		return "", false
	}
}

// streamAcceptMatch is the controlling Accept item found so far for one specificity tier (exact
// subtype, "type/*" wildcard, or bare "*/*" wildcard) in [matchStreamAccept]'s scan.
type streamAcceptMatch struct {
	value string
	found bool
	zero  bool
}

// consider records item as this tier's controlling reference if none has been recorded yet. Only the
// first occurrence within a tier is kept: RFC 9110 §12.5.1 lets the most specific tier present control
// regardless of order, and within one tier this package has no further precedence rule to break a tie
// between duplicate references, so the first one found is as good as any.
func (m *streamAcceptMatch) consider(zero bool, value string) {
	if !m.found {
		m.found, m.zero, m.value = true, zero, value
	}
}

// NewMediaFromContentType parses the request Content-Type header and resolves a streaming decoder.
//
// Unlike [Content.NewMediaFromContentType], an unregistered or unparsable media type returns
// [ErrUnsupportedMedia] instead of falling back to JSON.
//
// This is the streaming counterpart of [content.Content.NewFromRequestBody]: after NewMedia resolves the
// media type, it re-checks the resolved kind against the same request-decode policy that guards
// unary request bodies (see the decoder-bounds rule in the package documentation), so a streaming media
// type that maps to a banned kind is rejected here even though NewMedia itself has no opinion on
// that policy. NewFromAccept and [NewMedia] are not policed this way: only this
// constructor decodes an untrusted request body (see [Media.NewDecoder]'s other caller, the
// client's response-stream path, which must keep decoding every registered kind).
func NewMediaFromContentType(req *http.Request, sm *stream.Map) (Media, error) {
	m, err := NewMedia(req.Header.Get(content.TypeKey), sm)
	if err != nil {
		return Media{}, err
	}

	// Re-resolve the kind rather than carrying it on Media, keeping the struct free of a field
	// that exists only for this one caller. A missing entry, a nil decoder constructor, or a banned kind
	// are all treated as unsupported so this fails closed: neither the missing-entry nor the banned-kind
	// arm can currently fire, because streamKinds' only mapping ("application/x-ndjson" to "json") already passed the
	// lookup inside NewMedia and "json" is not a banned kind. Both stay as guards for a future
	// streamKinds entry that maps to a banned kind, or a NewMedia change that admits an unknown media
	// type.
	kind, ok := streamKinds[m.String()]
	if !ok || m.NewDecoder == nil || !policy.CanDecode(kind) {
		return Media{}, ErrUnsupportedMedia
	}

	return m, nil
}
