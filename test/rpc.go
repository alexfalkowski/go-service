package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nh "github.com/alexfalkowski/go-service/net/http"
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

func (*SuccessHandler) Handle(ctx nh.Context, r *Request) (*Response, error) {
	name := ctx.Request().URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Greeting: &s}, nil
}

type ProtobufHandler struct{}

func (*ProtobufHandler) Handle(_ nh.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func (*ProtobufHandler) Error(_ nh.Context, err error) *v1.SayHelloResponse {
	return &v1.SayHelloResponse{Message: err.Error()}
}

type ErrorHandler struct{}

func (*ErrorHandler) Handle(_ nh.Context, _ *Request) (*Response, error) {
	return nil, nh.Error(http.StatusServiceUnavailable, "ohh no")
}
