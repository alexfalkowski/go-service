package content

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrUnsupportedRequestMedia is returned when a request body uses a media type
// that is intentionally not decoded from public HTTP requests.
var ErrUnsupportedRequestMedia = errors.New("content: unsupported request media")

// ErrUnsupportedStreamMedia is returned when a media type has no streaming kind registered.
//
// Unlike single-value [Media] resolution, stream media resolution never falls back to JSON: silently
// degrading a streaming route to a different wire format would hide the mismatch from the caller.
var ErrUnsupportedStreamMedia = errors.New("content: unsupported stream media")

// ErrBidiRequiresHTTP2 is returned when a bidirectional streaming route is called over HTTP/1.x.
//
// HTTP/1.1 request bodies are buffered by intermediaries and the Go client transport before the
// response starts, so a handler that both reads the request stream and writes the response stream
// hangs rather than failing outright. Bidirectional streaming routes require HTTP/2 (including h2c).
var ErrBidiRequiresHTTP2 = errors.New("content: bidirectional streaming requires HTTP/2")
