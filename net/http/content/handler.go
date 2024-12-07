package content

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/errors"
	nh "github.com/alexfalkowski/go-service/net/http"
	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	"github.com/alexfalkowski/go-service/runtime"
)

// Handler for content.
type Handler func(ctx context.Context) (any, error)

// NewHandler for content.
func (c *Content) NewHandler(prefix string, handler Handler) func(res http.ResponseWriter, req *http.Request) {
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

		ct := c.NewFromRequest(req)

		ctx = hc.WithEncoder(ctx, ct.Encoder)
		res.Header().Add(TypeKey, ct.Media)

		data, err := handler(ctx)
		runtime.Must(err)

		if data == nil {
			return
		}

		err = ct.Encoder.Encode(res, data)
		runtime.Must(err)
	}

	return h
}
