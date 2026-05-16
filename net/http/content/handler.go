package content

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/ptr"
)

// RequestHandler handles a decoded request and returns a response model.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler builds a handler that decodes the request body and encodes the response.
//
// Context population:
// The handler attaches request-scoped values to the context via net/http/meta:
//   - the original *http.Request,
//   - the http.ResponseWriter, and
//   - the selected encoder.
//
// Content negotiation:
// The encoder is selected based on the request Content-Type, falling back to Accept when Content-Type is absent.
// The response Content-Type header is set to the negotiated media type.
//
// Errors:
// If request decoding fails, NewRequestHandler converts the decode error into a 400 Bad Request using
// net/http/status, allowing the response to be rendered consistently by status.WriteError.
//
// Successful responses are encoded into a pooled in-memory buffer before being written to the live
// response writer, so encode failures do not leak partial success bodies.
func NewRequestHandler[Req any, Res any](cont *Content, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		request := meta.Request(ctx)
		mediaType := cont.NewFromContentType(request)
		if err := mediaType.Encoder.Decode(request.Body, req); err != nil {
			return nil, status.BadRequestError(err)
		}

		return handler(ctx, req)
	})
}

// Handler handles a request without a request body and returns a response model.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler builds a handler that encodes the response and writes errors using status helpers.
//
// Context population and response content negotiation are the same as NewRequestHandler (see its documentation).
//
// Successful responses are encoded into a pooled in-memory buffer before being written to the live
// response writer, so encode failures do not leak partial success bodies.
func NewHandler[Res any](cont *Content, handler Handler[Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

// NotFoundHandler returns a not-found handler that writes the standard content error response.
func NotFoundHandler() http.NotFoundHandler {
	return func(res http.ResponseWriter, _ *http.Request) bool {
		err := status.Error(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		_ = status.WriteError(res, err)

		return true
	}
}

func newHandler[Res any](cont *Content, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		mediaType := cont.NewFromRequest(req)
		ctx = meta.WithContent(ctx, req, res, mediaType.Encoder)
		res.Header().Add(TypeKey, media.WithUTF8(mediaType.Type))

		data, err := handler(ctx)
		if err != nil {
			_ = status.WriteError(res, err)

			return
		}

		buffer := cont.pool.Get()
		defer cont.pool.Put(buffer)

		if err := mediaType.Encoder.Encode(buffer, data); err != nil {
			_ = status.WriteError(res, err)

			return
		}

		_, _ = buffer.WriteTo(res)
	}
}
