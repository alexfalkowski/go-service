package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/mime"
	"github.com/alexfalkowski/go-service/v2/net/http"
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
//   - gRPC status errors mapped to HTTP codes.
//
// Write behavior:
// The error message is written as a single line (via fmt.Fprintln) containing err.Error().
// If writing the body fails, WriteError does not attempt to write a secondary error response.
//
// This helper is conceptually similar to net/http.Error but uses go-service status code extraction
// and the dedicated error media type.
func WriteError(ctx context.Context, res http.ResponseWriter, err error) {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", mime.ErrorMediaType)
	header.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(Code(err))

	if _, err := fmt.Fprintln(res, err.Error()); err != nil {
		meta.WithAttribute(ctx, "writeError", meta.Error(err))
	}
}

// WriteText writes a plain-text success response to res.
//
// It clears any precomputed Content-Length, sets "Content-Type" to mime.TextMediaType,
// sets "X-Content-Type-Options: nosniff", writes HTTP 200 OK, and emits text followed by
// a trailing newline via fmt.Fprintln.
//
// If writing the body fails, WriteText does not attempt to write a secondary response.
func WriteText(ctx context.Context, res http.ResponseWriter, text string) {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", mime.TextMediaType)
	header.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintln(res, text); err != nil {
		meta.WithAttribute(ctx, "writeError", meta.Error(err))
	}
}
