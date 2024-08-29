package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/encoding"
	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

// Handler for content.
type Handler func(ctx context.Context) any

// NewHandler for content.
func NewHandler(prefix string, enc *encoding.Map, handler Handler) func(res http.ResponseWriter, req *http.Request) {
	h := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = hc.WithRequest(ctx, req)
		ctx = hc.WithResponse(ctx, res)

		defer func() {
			if r := recover(); r != nil {
				err := errors.Prefix(prefix, runtime.ConvertRecover(r))
				nh.WriteError(ctx, res, err, status.Code(err))
			}
		}()

		ct := NewFromRequest(req, enc)

		ctx = hc.WithEncoder(ctx, ct.Encoder)
		res.Header().Add(TypeKey, ct.Media)

		data := handler(ctx)

		err := ct.Encoder.Encode(res, data)
		runtime.Must(err)
	}

	return h
}
