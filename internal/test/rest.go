package test

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/client"
	"github.com/alexfalkowski/go-service/v2/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
)

// RegisterHandlers registers DELETE and GET REST handlers for the service-prefixed path.
func (w *World) RegisterHandlers[Res any](path string, h unary.Handler[Res]) {
	w.RestServer.Delete(http.Pattern(Name, path), h)
	w.RestServer.Get(http.Pattern(Name, path), h)
}

// RegisterRequestHandlers registers POST, PUT, and PATCH REST handlers for the service-prefixed path.
func (w *World) RegisterRequestHandlers[Req any, Res any](path string, h unary.RequestHandler[Req, Res]) {
	w.RestServer.Post(http.Pattern(Name, path), h)
	w.RestServer.Put(http.Pattern(Name, path), h)
	w.RestServer.Patch(http.Pattern(Name, path), h)
}

// RestInvalidStatusCode writes an internal server error directly to the response and returns no payload.
func RestInvalidStatusCode(ctx context.Context) (*Response, error) {
	res := meta.Response(ctx)
	res.WriteHeader(http.StatusInternalServerError)

	return nil, nil
}

// RestNoContent returns no body and no error so callers can exercise empty-success responses.
func RestNoContent(_ context.Context) (*Response, error) {
	return nil, nil
}

// RestRequestInvalidStatusCode writes an internal server error directly to the response for request-body handlers.
func RestRequestInvalidStatusCode(ctx context.Context, _ *Request) (*Response, error) {
	res := meta.Response(ctx)
	res.WriteHeader(http.StatusInternalServerError)

	return nil, nil
}

// RestRequestNoContent returns no body and no error for request-body handlers.
func RestRequestNoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// RestContent builds a greeting from the `name` query parameter and echoes camel-cased request metadata.
func RestContent(ctx context.Context) (*Response, error) {
	req := meta.Request(ctx)
	_ = meta.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, meta.NoPrefix), Greeting: s}, nil
}

// RestRequestContent builds a greeting from the request body and echoes camel-cased request metadata.
func RestRequestContent(ctx context.Context, req *Request) (*Response, error) {
	name := cmp.Or(req.Name, "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, meta.NoPrefix), Greeting: s}, nil
}

// RestRequestProtobuf returns a protobuf greeting response for REST-to-protobuf content tests.
func RestRequestProtobuf(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	name := cmp.Or(r.GetName(), "Bob")
	s := "Hello " + name

	return &v1.SayHelloResponse{Message: s}, nil
}

// RestError returns ErrInvalid so REST tests can exercise mapped error responses.
func RestError(_ context.Context) (*Response, error) {
	return nil, ErrInvalid
}

// RestRequestError returns ErrInvalid for request-body REST handlers.
func RestRequestError(_ context.Context, _ *Request) (*Response, error) {
	return nil, ErrInvalid
}

func restClient(httpClient *http.Client, os *options) *rest.Client {
	opts := []client.ClientOption{client.WithRedirect(client.RedirectIgnore)}
	if os.rest {
		opts = append(opts, client.WithRoundTripper(httpClient.Transport))
	}

	return rest.NewClient(client.NewClient(UnaryContent, StreamContent, Pool, opts...))
}
