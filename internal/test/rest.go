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

// WithWorldRest for test.
func WithWorldRest() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.rest = true
	})
}

// RegisterHandlers for test.
func RegisterHandlers[Res any](path string, h content.Handler[Res]) {
	rest.Delete(http.Pattern(Name, path), h)
	rest.Get(http.Pattern(Name, path), h)
}

// RegisterRequestHandlers for test.
func RegisterRequestHandlers[Req any, Res any](path string, h content.RequestHandler[Req, Res]) {
	rest.Post(http.Pattern(Name, path), h)
	rest.Put(http.Pattern(Name, path), h)
	rest.Patch(http.Pattern(Name, path), h)
}

// RestInvalidStatusCode for test.
func RestInvalidStatusCode(ctx context.Context) (*Response, error) {
	res := meta.Response(ctx)
	res.WriteHeader(http.StatusInternalServerError)

	return nil, nil
}

// RestNoContent for test.
func RestNoContent(_ context.Context) (*Response, error) {
	return nil, nil
}

// RestRequestInvalidStatusCode for test.
func RestRequestInvalidStatusCode(ctx context.Context, _ *Request) (*Response, error) {
	res := meta.Response(ctx)
	res.WriteHeader(http.StatusInternalServerError)

	return nil, nil
}

// RestRequestNoContent for test.
func RestRequestNoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// RestContent for test.
func RestContent(ctx context.Context) (*Response, error) {
	req := meta.Request(ctx)
	_ = meta.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, meta.NoPrefix), Greeting: s}, nil
}

// RestRequestContent for test.
func RestRequestContent(ctx context.Context, req *Request) (*Response, error) {
	name := cmp.Or(req.Name, "Bob")
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, meta.NoPrefix), Greeting: s}, nil
}

// RestRequestProtobuf for test.
func RestRequestProtobuf(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	name := cmp.Or(r.GetName(), "Bob")
	s := "Hello " + name

	return &v1.SayHelloResponse{Message: s}, nil
}

// RestError for test.
func RestError(_ context.Context) (*Response, error) {
	return nil, ErrInvalid
}

// RestRequestError for test.
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

func restClient(client *Client, os *worldOpts) *rest.Client {
	if os.rest {
		return rest.NewClient(
			rest.WithClientRoundTripper(client.NewHTTP().Transport),
			rest.WithClientTimeout("10s"),
		)
	}

	return rest.NewClient()
}
