package content

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/types/ptr"
)

// RequestHandler handles a decoded request and returns a response model.
type RequestHandler[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// NewRequestHandler builds a handler that decodes the request body and encodes the response.
//
// Context population:
// The handler attaches request-scoped values to the context via net/http/meta:
//   - the original *http.Request,
//   - the http.ResponseWriter, and
//   - the selected encoder (derived from the request Content-Type).
//
// Content negotiation:
// The encoder is selected based on the request Content-Type (via (*Content).NewFromRequest).
// The response Content-Type header is set to the negotiated media type.
//
// Errors:
// If request decoding fails, NewRequestHandler converts the decode error into a 400 Bad Request using
// net/http/status, allowing the response to be rendered consistently by status.WriteError.
//
// Constraints:
// The negotiated media type must resolve to a non-nil encoder. If the request Content-Type resolves to
// the error media subtype ("error"), Media.Encoder will be nil and decoding will panic. In practice,
// clients should never send Content-Type "text/error"; that media type is reserved for error responses.
func NewRequestHandler[Req any, Res any](cont *Content, handler RequestHandler[Req, Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		req := ptr.Zero[Req]()

		encoder := meta.Encoder(ctx)
		request := meta.Request(ctx)

		if err := encoder.Decode(request.Body, req); err != nil {
			return nil, status.BadRequestError(err)
		}

		return handler(ctx, req)
	})
}

// Handler handles a request without a request body and returns a response model.
type Handler[Res any] func(ctx context.Context) (*Res, error)

// NewHandler builds a handler that encodes the response and writes errors using status helpers.
//
// Context population and content negotiation are the same as NewRequestHandler (see its documentation).
func NewHandler[Res any](cont *Content, handler Handler[Res]) http.HandlerFunc {
	return newHandler(cont, func(ctx context.Context) (*Res, error) {
		return handler(ctx)
	})
}

func newHandler[Res any](cont *Content, handler func(ctx context.Context) (*Res, error)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = meta.WithRequest(ctx, req)
		ctx = meta.WithResponse(ctx, res)
		media := cont.NewFromRequest(req)

		ctx = meta.WithEncoder(ctx, media.Encoder)
		res.Header().Add(TypeKey, media.Type)

		data, err := handler(ctx)
		if err != nil {
			status.WriteError(ctx, res, err)

			return
		}

		if err := media.Encoder.Encode(res, data); err != nil {
			status.WriteError(ctx, res, err)

			return
		}
	}
}
