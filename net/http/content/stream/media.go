package stream

import (
	"mime"

	"github.com/alexfalkowski/go-service/v2/encoding/stream"
	"github.com/alexfalkowski/go-service/v2/net/http"
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
// Streaming currently supports one such type. Keeping the direct lookup separate leaves
// matchStreamAccept responsible only for full Accept-list negotiation.
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

const (
	noStreamMatch = iota
	bareStreamMatch
	majorStreamMatch
	exactStreamMatch
)

// streamAcceptMatch is a matching streaming media range and its RFC 9110 specificity.
type streamAcceptMatch struct {
	mediaType   string
	specificity int
	zeroQuality bool
}

// matchStreamAccept reports whether header, an Accept header value, is satisfiable for a registered
// streamKinds entry, and if so, which media type to resolve.
//
// It finds the most specific reference to the parameterless producible type present anywhere in the list —
// an exact subtype match without non-quality parameters, else a "type/*" wildcard [http.IsAcceptWildcard]
// reports as satisfied by [ndjsonType] without non-quality parameters ("text/*" does not satisfy an
// "application/x-ndjson" route the way "*/*" or "application/*" does), else the bare "*/*" wildcard —
// and returns satisfiable only if that one reference is not [http.IsAcceptZeroQuality]; a less specific
// reference's own quality, zero or not, is irrelevant once a more specific reference is found. When two
// references tie on specificity, the zero-quality one controls, so an explicit exclusion is not masked by
// a duplicate positive entry for the same range. This matches RFC 9110 §12.5.1's "most specific reference"
// rule regardless of list order. When satisfiable, it returns
// the exact match's media type if the controlling reference was one, otherwise [media.NDJSON]. An
// unparsable item is skipped rather than rejecting the whole list, the same way an unparsable single Accept
// value already falls through to [Content.NewFromMedia]'s own rejection when nothing else in the list matches.
func matchStreamAccept(header string) (string, bool) {
	var best streamAcceptMatch

	for _, item := range http.AcceptItems(header) {
		candidate := matchStreamRange(item)
		if candidate.specificity > best.specificity || (candidate.specificity == best.specificity && candidate.zeroQuality) {
			best = candidate
		}
	}

	return best.mediaType, best.specificity != noStreamMatch && !best.zeroQuality
}

func matchStreamRange(item string) streamAcceptMatch {
	value, ok := parseParameterlessAcceptRange(item)
	if !ok {
		return streamAcceptMatch{}
	}

	match := streamAcceptMatch{zeroQuality: http.IsAcceptZeroQuality(item)}
	if _, ok := streamKinds[value.String()]; ok {
		match.mediaType = value.String()
		match.specificity = exactStreamMatch

		return match
	}

	if !http.IsAcceptWildcard(value, ndjsonType) {
		return streamAcceptMatch{}
	}

	match.mediaType = media.NDJSON
	if value.Major() == "*" {
		match.specificity = bareStreamMatch
	} else {
		match.specificity = majorStreamMatch
	}

	return match
}

func parseParameterlessAcceptRange(item string) (media.Type, bool) {
	value, err := media.Parse(item)
	if err != nil || hasNonQualityParameters(item) {
		return media.Type{}, false
	}

	return value, true
}

// hasNonQualityParameters reports whether item requires media parameters the canonical parameterless
// representation does not provide.
func hasNonQualityParameters(item string) bool {
	_, parameters, err := mime.ParseMediaType(item)
	if err != nil {
		return false
	}

	delete(parameters, "q")

	return len(parameters) != 0
}
