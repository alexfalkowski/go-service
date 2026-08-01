// Package accept provides RFC 9110 Accept header parsing and matching helpers used by go-service.
//
// [Items] and [First] split an HTTP list-valued header value, such as Accept, into its comma-separated
// items, honoring the quoted-comma rule from RFC 9110 §5.6.1. [IsZeroQuality] and [IsWildcard] answer the
// two questions a caller needs to decide whether one Accept item is satisfiable for a producible media
// type: whether the item explicitly excludes itself via a zero quality value (RFC 9110 §12.4.2), and
// whether the item is a wildcard media range that matches the producible type's major type (RFC 9110
// §12.5.1).
//
// This package only answers those two questions; it does not rank items by quality value or wildcard
// precedence. A caller that needs to decide whether an Accept list is satisfiable for the media type(s)
// it can produce combines these helpers itself — see
// [github.com/alexfalkowski/go-service/v2/net/http/content/stream.NewFromAccept] for an example.
package accept
