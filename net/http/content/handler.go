package content

import (
	"context"
	"net/http"
	"net/url"

	hc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
)

type (
	handler func(ctx context.Context) (any, error)

	// Handler is a handler with a generic request and response.
	Handler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)
)

// NewQueryHandler for content.
func NewQueryHandler[Req any, Res any](cont *Content, prefix string, handler Handler[Req, Res]) http.HandlerFunc {
	return cont.handler(prefix, func(ctx context.Context) (any, error) {
		e := hc.Encoder(ctx)
		req := hc.Request(ctx)

		b := pool.Get()
		defer pool.Put(b)

		q := convertQuery(req.URL.Query())

		if err := e.Encode(b, q); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		var rq Req
		ptr := &rq

		if err := e.Decode(b, ptr); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		rs, err := handler(ctx, ptr)
		if err != nil {
			return nil, err
		}

		return rs, nil
	})
}

// NewBodyHandler for content.
func NewBodyHandler[Req any, Res any](cont *Content, prefix string, handler Handler[Req, Res]) http.HandlerFunc {
	return cont.handler(prefix, func(ctx context.Context) (any, error) {
		var rq Req
		ptr := &rq

		e := hc.Encoder(ctx)
		req := hc.Request(ctx)

		if err := e.Decode(req.Body, ptr); err != nil {
			return nil, status.Error(http.StatusBadRequest, err.Error())
		}

		rs, err := handler(ctx, ptr)
		if err != nil {
			return nil, err
		}

		return rs, nil
	})
}

func convertQuery(query url.Values) map[string]any {
	values := make(map[string]any, len(query))

	for k, v := range query {
		if len(v) == 1 {
			values[k] = v[0]
		} else {
			values[k] = v
		}
	}

	return values
}
