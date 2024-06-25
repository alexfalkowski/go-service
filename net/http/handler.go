package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/meta"
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
	Handle(ctx context.Context, req *Req) (*Res, error)

	// Error for this handler.
	Error(ctx context.Context, err error) *Res

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
func Handle[Req any, Res any](pattern string, handler Handler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		c, k := kind(req)
		m := mar.Get(k)

		res.Header().Add("Content-Type", c)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "readAllError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrReadAll, err), handler)

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			ctx = meta.WithAttribute(ctx, "unmarshalError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrUnmarshal, err), handler)

			return
		}

		rs, err := handler.Handle(ctx, ptr)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "handleError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrHandle, err), handler)

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "marshalError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrMarshal, err), handler)

			return
		}

		res.Write(d)
	}

	mux.HandleFunc(pattern, h)
}

func kind(req *http.Request) (string, string) {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func writeError[Req any, Res any](ctx context.Context, res http.ResponseWriter, m marshaller.Marshaller, err error, h Handler[Req, Res]) {
	res.WriteHeader(h.Status(err))

	d, err := m.Marshal(h.Error(ctx, err))
	if err != nil {
		meta.WithAttribute(ctx, "writeError", meta.Error(err))

		return
	}

	res.Write(d)
}
