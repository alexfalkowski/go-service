package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/content"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

var (
	mux *http.ServeMux
	enc *encoding.Map
)

// Register for rpc.
func Register(mu *http.ServeMux, en *encoding.Map) {
	mux, enc = mu, en
}

// UnaryHandler for rpc.
type UnaryHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// Unary for rpc.
func Unary[Req any, Res any](path string, handler UnaryHandler[Req, Res]) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix("rpc", runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := content.NewFromRequest(req)

		m, err := ct.Marshaller(enc)
		runtime.Must(err)

		res.Header().Add(content.TypeKey, ct.Media)

		body, err := io.ReadAll(req.Body)
		runtime.Must(err)

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var rq Req
		ptr := &rq

		err = m.Unmarshal(body, ptr)
		runtime.Must(err)

		rs, err := handler(ctx, ptr)
		runtime.Must(err)

		d, err := m.Marshal(rs)
		runtime.Must(err)

		res.Write(d)
	}

	mux.HandleFunc("POST "+path, h)
}
