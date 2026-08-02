package stream

import (
	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http/accept"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
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

// knownMedia resolves an exact built-in streaming media type without parsing its header value.
//
// Streaming currently supports one such type, but keeping this separate from matchStreamAccept
// makes the concrete-header fast path explicit and leaves the latter responsible only for full
// Accept-list negotiation.
func knownMedia(mediaType string, enc *stream.Map) (Media, bool) {
	if mediaType != media.NDJSON {
		return Media{}, false
	}

	codec := enc.Get(streamKinds[media.NDJSON])
	return Media{NewEncoder: codec.Encoder, NewDecoder: codec.Decoder, Type: ndjsonType}, true
}

func usableMedia(resolved Media) (Media, error) {
	if resolved.NewEncoder == nil && resolved.NewDecoder == nil {
		return Media{}, ErrUnsupportedMedia
	}

	return resolved, nil
}

func encoderMedia(resolved Media, err error) (Media, error) {
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
// [Content.NewFromMedia]'s own rejection when nothing else in the list matches.
func matchStreamAccept(header string) (string, bool) {
	var exact, major, bare mediaRangeMatch

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

// mediaRangeMatch is the controlling media range found so far for one specificity tier (exact subtype,
// "type/*" wildcard, or bare "*/*" wildcard) in [matchStreamAccept]'s scan.
type mediaRangeMatch struct {
	value string
	found bool
	zero  bool
}

// consider records item as this tier's controlling reference if none has been recorded yet. Only the
// first occurrence within a tier is kept: RFC 9110 §12.5.1 lets the most specific tier present control
// regardless of order, and within one tier this package has no further precedence rule to break a tie
// between duplicate references, so the first one found is as good as any.
func (m *mediaRangeMatch) consider(zero bool, value string) {
	if !m.found {
		m.found, m.zero, m.value = true, zero, value
	}
}
