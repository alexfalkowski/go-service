package http

import (
	"bytes"
	"context"
	"errors"
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

// Errorer for HTTP.
type Errorer[Res any] interface {
	// Error for this handler.
	Error(ctx context.Context, err error) *Res

	// Status code from error.
	Status(err error) int
}

// NewHandler for HTTP.
func NewHandler[Req any, Res any](mux ServeMux, mar *marshaller.Map, err Errorer[Res]) *Handler[Req, Res] {
	return &Handler[Req, Res]{mux: mux, mar: mar, err: err}
}

// Handler for HTTP.
type Handler[Req any, Res any] struct {
	mux ServeMux
	mar *marshaller.Map
	err Errorer[Res]
}

// Handle func for request/response.
type Handle[Req any, Res any] func(context.Context, *Req) (*Res, error)

// Handle for HTTP.
func (s *Handler[Req, Res]) Handle(verb, pattern string, fn Handle[Req, Res]) error {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		c, k := s.kind(req)
		m := s.mar.Get(k)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "readAllError", meta.Error(err))
			s.error(ctx, res, m, ErrReadAll)

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			ctx = meta.WithAttribute(ctx, "unmarshalError", meta.Error(err))
			s.error(ctx, res, m, ErrUnmarshal)

			return
		}

		res.Header().Add("Content-Type", c)

		rs, err := fn(ctx, ptr)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "handleError", meta.Error(err))
			s.error(ctx, res, m, ErrHandle)

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			ctx = meta.WithAttribute(ctx, "marshalError", meta.Error(err))
			s.error(ctx, res, m, ErrMarshal)

			return
		}

		res.Write(d)
	}

	return s.mux.Handle(verb, pattern, h)
}

func (s *Handler[Req, Res]) kind(req *http.Request) (string, string) {
	t, err := ct.GetMediaType(req)
	if err != nil {
		return "application/json", "json"
	}

	return t.String(), t.Subtype
}

func (s *Handler[Req, Res]) error(ctx context.Context, res http.ResponseWriter, m marshaller.Marshaller, err error) {
	d, _ := m.Marshal(s.err.Error(ctx, err))

	res.WriteHeader(s.err.Status(err))
	res.Write(d)
}
