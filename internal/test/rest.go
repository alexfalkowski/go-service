package test

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
)

// RegisterHandlers registers DELETE and GET REST handlers for the service-prefixed path.
func RegisterHandlers[Res any](path string, h content.Handler[Res]) {
	rest.Delete(http.Pattern(Name, path), h)
	rest.Get(http.Pattern(Name, path), h)
}

// RegisterRequestHandlers registers POST, PUT, and PATCH REST handlers for the service-prefixed path.
func RegisterRequestHandlers[Req any, Res any](path string, h content.RequestHandler[Req, Res]) {
	rest.Post(http.Pattern(Name, path), h)
	rest.Put(http.Pattern(Name, path), h)
	rest.Patch(http.Pattern(Name, path), h)
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

func registerRest(mux *http.ServeMux) {
	rest.Register(rest.RegisterParams{
		Mux:     mux,
		Pool:    Pool,
		Content: Content,
	})
}

func restClient(client *http.Client, os *worldOpts) *rest.Client {
	if os.rest {
		return rest.NewClient(
			rest.WithClientRoundTripper(client.Transport),
			rest.WithClientTimeout("10s"),
		)
	}

	return rest.NewClient()
}
