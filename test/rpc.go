package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
)

type Request struct {
	Name string
}

type Response struct {
	Meta     meta.Map
	Greeting *string
}

type SuccessHandler struct{}

func (*SuccessHandler) Handle(ctx rpc.Context, r *Request) (*Response, error) {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Greeting: &s}, nil
}

type ProtobufHandler struct{}

func (*ProtobufHandler) Handle(_ rpc.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func (*ProtobufHandler) Error(_ rpc.Context, err error) *v1.SayHelloResponse {
	return &v1.SayHelloResponse{Message: err.Error()}
}

type ErrorHandler struct{}

func (*ErrorHandler) Handle(_ rpc.Context, _ *Request) (*Response, error) {
	return nil, rpc.Error(http.StatusServiceUnavailable, "ohh no")
}
