package rpc

import (
	"bytes"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/marshaller"
	"github.com/alexfalkowski/go-service/net/http/content"
)

var (
	mux *http.ServeMux
	mar *marshaller.Map
)

// Register for HTTP.
func Register(mu *http.ServeMux, ma *marshaller.Map) {
	mux, mar = mu, ma
}

// Handler for HTTP.
type Handler[Req any, Res any] interface {
	// Handle the request/response.
	Handle(ctx Context, req *Req) (*Res, error)
}

// Handler for HTTP.
func Handle[Req any, Res any](path string, handler Handler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req.Context(), req, res)
		ct := content.NewFromRequest(req)

		m, err := ct.Marshaller(mar)
		if err != nil {
			WriteError(ctx, errors.Prefix("rpc marshaller", err))

			return
		}

		res.Header().Add(content.TypeKey, ct.Media)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			WriteError(ctx, errors.Prefix("rpc read", err))

			return
		}

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		if err := m.Unmarshal(body, ptr); err != nil {
			WriteError(ctx, errors.Prefix("rpc unmarshal", err))

			return
		}

		rs, err := handler.Handle(ctx, ptr)
		if err != nil {
			WriteError(ctx, errors.Prefix("rpc handle", err))

			return
		}

		d, err := m.Marshal(rs)
		if err != nil {
			WriteError(ctx, errors.Prefix("rpc marshal", err))

			return
		}

		res.Write(d)
	}

	mux.HandleFunc("POST "+path, h)
}
