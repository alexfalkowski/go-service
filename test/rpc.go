package test

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/status"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
)

type Request struct {
	Name string
}

type Response struct {
	Meta     meta.Map
	Greeting *string
}

func SuccessSayHello(ctx context.Context, r *Request) (*Response, error) {
	req := nc.Request(ctx)

	name := req.URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: &s}, nil
}

func ProtobufSayHello(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func ErrorSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, status.Error(http.StatusServiceUnavailable, "ohh no")
}
