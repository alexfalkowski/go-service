package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/marshaller"
	ct "github.com/elnormous/contenttype"
)

var (
	// ErrReadAll for HTTP.
	ErrReadAll = errors.New("invalid body read")

	// ErrUnmarshal for HTTP.
	ErrUnmarshal = errors.New("invalid unmarshal")

	// ErrHandle for HTTP.
	ErrHandle = errors.New("invalid handle")

	// ErrMarshal for HTTP.
	ErrMarshal = errors.New("invalid marshal")

	contentTypeKey = "Content-Type"
)

// Handler for HTTP.
type Handler[Req any, Res any] interface {
	// Handle func for request/response.
	Handle(ctx Context, req *Req) (*Res, error)
}

var (
	mux *http.ServeMux
	mar *marshaller.Map
)

// Register for HTTP.
func Register(mu *http.ServeMux, ma *marshaller.Map) {
	mux, mar = mu, ma
}

// Handler for HTTP.
func Handle[Req any, Res any](path string, handler Handler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := newContext(req.Context(), req, res)
		c, k := kindFromRequest(req)
		m := mar.Get(k)

		res.Header().Add(contentTypeKey, c)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			writeError(ctx, fmt.Errorf("%w: %w", ErrReadAll, err))

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			writeError(ctx, fmt.Errorf("%w: %w", ErrUnmarshal, err))

			return
		}

		rs, err := handler.Handle(ctx, ptr)
		if err != nil {
			writeError(ctx, fmt.Errorf("%w: %w", ErrHandle, err))

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			writeError(ctx, fmt.Errorf("%w: %w", ErrMarshal, err))

			return
		}

		res.Write(d)
	}

	mux.HandleFunc("POST "+path, h)
}

func kindFromRequest(req *http.Request) (string, string) {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func writeError(ctx Context, err error) {
	http.Error(ctx.Response(), err.Error(), Code(err))
}
