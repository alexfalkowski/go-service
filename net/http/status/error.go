package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
)

var (
	errorContentType = media.MustParse(media.Error).WithUTF8()
	textContentType  = media.MustParse(media.Text).WithUTF8()
)

// WriteError records err for request-scoped operator diagnostics, then writes its safe response to res.
//
// The first error written through ctx is retained because it determines the client response. The response
// rendering and error-return behavior are otherwise identical to previous releases of WriteError. It writes a
// plain-text `text/error; charset=utf-8` response with `X-Content-Type-Options: nosniff`. The status code is
// derived with [Code]. The body contains the first safe message in err's chain, or the default message for that
// code. Use [SafeError] to retain an internal cause without exposing it to the client. WriteError returns a
// body-write failure without attempting a secondary response.
func WriteError(ctx context.Context, res http.ResponseWriter, err error) error {
	if state, _ := ctx.Value(requestErrorKey).(*requestError); state != nil && state.err == nil {
		state.err = err
	}

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
