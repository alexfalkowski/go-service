package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	"github.com/alexfalkowski/go-service/net/http/content"
	hc "github.com/alexfalkowski/go-service/net/http/context"
)

// UnaryHandler for rpc.
type UnaryHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// Unary for rpc.
func Unary[Req any, Res any](path string, handler UnaryHandler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)
		ct := content.NewFromRequest(req)

		m, err := ct.Marshaller(enc)
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

		rs, err := handler(ctx, ptr)
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
