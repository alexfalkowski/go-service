package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

var (
	errorContentType = media.MustParse(media.Error).WithUTF8()
	textContentType  = media.MustParse(media.Text).WithUTF8()
)

// WriteError writes an error response to res.
//
// Content-Type:
// WriteError always writes the response as a plain-text error payload using the go-service specific
// error media type "text/error; charset=utf-8" and sets "X-Content-Type-Options: nosniff".
//
// Status code selection:
// The HTTP status code is derived from err using Code(err), which understands:
//   - errors implementing Coder,
//   - errors created by this package, and
//   - raw *[http.MaxBytesError] values mapped to 413 Request Entity Too Large, and
//   - gRPC status errors mapped to HTTP codes.
//
// Write behavior:
// The error message is written as a single line (via [fmt.Fprintln]) containing the first SafeMessage in
// err's chain. If no safe message is available, WriteError uses the default safe HTTP status message for Code(err).
// Use SafeError to preserve an internal cause while returning the default safe HTTP status message to the client.
// If writing the body fails, WriteError returns the write error and does not attempt to write a
// secondary error response.
//
// This helper is conceptually similar to [net/http.Error] but uses go-service status code extraction
// and the dedicated error media type.
func WriteError(res http.ResponseWriter, err error) error {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", errorContentType)
	header.Set("X-Content-Type-Options", "nosniff")

	code := Code(err)
	res.WriteHeader(code)

	_, writeErr := fmt.Fprintln(res, errors.SafeMessage(err, DefaultMessage(code)))
	return writeErr
}

// WriteText writes a plain-text success response to res.
//
// It clears any precomputed Content-Length, sets "Content-Type" to [media.Text] with UTF-8 charset,
// sets "X-Content-Type-Options: nosniff", writes HTTP 200 OK, and emits text followed by
// a trailing newline via [fmt.Fprintln].
//
// If writing the body fails, WriteText returns the write error and does not attempt to write a
// secondary response.
func WriteText(res http.ResponseWriter, text string) error {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", textContentType)
	header.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(http.StatusOK)

	_, err := fmt.Fprintln(res, text)
	return err
}
