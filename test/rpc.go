package test

import (
	"cmp"
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
)

// Request for test.
type Request struct {
	Name string
}

// Response for test.
type Response struct {
	Meta     meta.Map
	Greeting string
}

// NoContent for test.
//
//nolint:nilnil
func NoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// SuccessSayHello for test.
func SuccessSayHello(ctx context.Context, r *Request) (*Response, error) {
	req := nc.Request(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), r.Name)
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: s}, nil
}

// ProtobufSayHello for test.
func ProtobufSayHello(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

// ErrorSayHello for test.
func ErrorSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, status.Error(http.StatusServiceUnavailable, "ohh no")
}
