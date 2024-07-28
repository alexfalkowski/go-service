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

type (
	// Router for rpc.
	Router struct {
		mux *http.ServeMux
		enc *encoding.Map
	}

	// Request for rpc.
	Request struct {
		mar  encoding.Marshaller
		data []byte
	}

	// Response for rpc.
	Response any

	// Handler for rpc.
	Handler func(ctx context.Context, req *Request) (Response, error)
)

// NewRouter for rpc.
func NewRouter(mux *http.ServeMux, enc *encoding.Map) *Router {
	return &Router{mux: mux, enc: enc}
}

// Unmarshal the request.
func (r *Request) Unmarshal(req any) {
	err := r.mar.Unmarshal(r.data, req)
	runtime.Must(err)
}

// Route for rpc.
func (r *Router) Route(path string, handler Handler) {
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

		m, err := ct.Marshaller(r.enc)
		runtime.Must(err)

		res.Header().Add(content.TypeKey, ct.Media)

		body, err := io.ReadAll(req.Body)
		runtime.Must(err)

		req.Body = io.NopCloser(bytes.NewBuffer(body))

		rq := &Request{data: body, mar: m}

		rs, err := handler(ctx, rq)
		runtime.Must(err)

		d, err := m.Marshal(rs)
		runtime.Must(err)

		res.Write(d)
	}

	r.mux.HandleFunc("POST "+path, h)
}
