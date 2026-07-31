package content

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrUnsupportedRequestMedia is returned when a request Content-Type cannot be decoded from a public
// HTTP request body. This covers three causes: the media type is intentionally denied (see the
// decoder-bounds rule in the package documentation), the media type does not parse, or the media type
// parses but resolves to no registered codec. See [Content.NewFromRequestBody].
var ErrUnsupportedRequestMedia = errors.New("content: unsupported request media")

// ErrUnsupportedStreamMedia is returned when a streaming media type has no streaming kind registered, or
// when [Content.NewStreamFromContentType] rejects a resolved kind under the same request-decode policy
// that guards [ErrUnsupportedRequestMedia].
//
// Unlike single-value [Media] resolution's outbound fallback, stream media resolution never falls back
// to JSON: silently degrading a streaming route to a different wire format would hide the mismatch from
// the caller. That contrast has narrowed on the inbound side too, now that an unknown or unregistered
// request Content-Type is rejected rather than decoded as JSON — see [Content.NewFromRequestBody].
var ErrUnsupportedStreamMedia = errors.New("content: unsupported stream media")

// ErrBidiRequiresHTTP2 is returned when a bidirectional streaming route is called over HTTP/1.x.
//
// HTTP/1.1 request bodies are buffered by intermediaries and the Go client transport before the
// response starts, so a handler that both reads the request stream and writes the response stream
// hangs rather than failing outright. Bidirectional streaming routes require HTTP/2 (including h2c).
var ErrBidiRequiresHTTP2 = errors.New("content: bidirectional streaming requires HTTP/2")
