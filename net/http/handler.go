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
)

// Handler for HTTP.
type Handler[Req any, Res any] interface {
	// Handle func for request/response.
	Handle(ctx Context, req *Req) (*Res, error)

	// Error for this handler.
	Error(ctx Context, err error) *Res

	// Status code from error.
	Status(err error) int
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
		c, k := kind(req)
		m := mar.Get(k)

		res.Header().Add("Content-Type", c)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			writeError(ctx, m, fmt.Errorf("%w: %w", ErrReadAll, err), handler)

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			writeError(ctx, m, fmt.Errorf("%w: %w", ErrUnmarshal, err), handler)

			return
		}

		rs, err := handler.Handle(ctx, ptr)
		if err != nil {
			writeError(ctx, m, fmt.Errorf("%w: %w", ErrHandle, err), handler)

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			writeError(ctx, m, fmt.Errorf("%w: %w", ErrMarshal, err), handler)

			return
		}

		res.Write(d)
	}

	mux.HandleFunc("POST "+path, h)
}

func kind(req *http.Request) (string, string) {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func writeError[Req any, Res any](ctx Context, m marshaller.Marshaller, err error, h Handler[Req, Res]) {
	res := ctx.Response()
	res.WriteHeader(h.Status(err))

	d, _ := m.Marshal(h.Error(ctx, err))
	res.Write(d)
}
