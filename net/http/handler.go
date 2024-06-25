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

type (
	// Errorer for HTTP.
	Errorer[Res any] interface {
		// Error for this handler.
		Error(ctx context.Context, err error) *Res

		// Status code from error.
		Status(err error) int
	}

	// Handle func for request/response.
	Handle[Req any, Res any] func(context.Context, *Req) (*Res, error)
)

var (
	mux *http.ServeMux
	mar *marshaller.Map
)

// RegisterHandler for HTTP.
func RegisterHandler(mu *http.ServeMux, ma *marshaller.Map) {
	mux, mar = mu, ma
}

// Handler for HTTP.
func Handler[Req any, Res any](pattern string, errorer Errorer[Res], fn Handle[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		c, k := kind(req)
		m := mar.Get(k)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "readAllError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrReadAll, err), errorer)

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			ctx = meta.WithAttribute(ctx, "unmarshalError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrUnmarshal, err), errorer)

			return
		}

		res.Header().Add("Content-Type", c)

		rs, err := fn(ctx, ptr)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "handleError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrHandle, err), errorer)

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "marshalError", meta.Error(err))
			writeError(ctx, res, m, fmt.Errorf("%w: %w", ErrMarshal, err), errorer)

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

func writeError[Res any](ctx context.Context, res http.ResponseWriter, m marshaller.Marshaller, err error, errorer Errorer[Res]) {
	d, _ := m.Marshal(errorer.Error(ctx, err))

	res.WriteHeader(errorer.Status(err))
	res.Write(d)
}
