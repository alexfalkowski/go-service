package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
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
// If writing the body fails, WriteError records the write failure as a meta attribute on ctx
// under the key "writeError". It does not attempt to write a secondary error response.
//
// This helper is conceptually similar to net/http.Error but uses go-service status code extraction
// and the dedicated error media type.
func WriteError(ctx context.Context, res http.ResponseWriter, err error) {
	header := res.Header()
	header.Del("Content-Length")
	header.Set("Content-Type", "text/error; charset=utf-8")
	header.Set("X-Content-Type-Options", "nosniff")

	res.WriteHeader(Code(err))

	if _, err := fmt.Fprintln(res, err.Error()); err != nil {
		meta.WithAttribute(ctx, "writeError", meta.Error(err))
	}
}
