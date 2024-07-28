package test

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
	"github.com/alexfalkowski/go-service/net/http/rpc"
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

func SuccessSayHello(ctx context.Context, request *rpc.Request) (rpc.Response, error) {
	var r Request

	request.Unmarshal(&r)

	req := nc.Request(ctx)

	name := req.URL.Query().Get("name")
	if name == "" {
		name = r.Name
	}

	s := "Hello " + name

	return &Response{Greeting: &s}, nil
}

func ProtobufSayHello(_ context.Context, request *rpc.Request) (rpc.Response, error) {
	var r v1.SayHelloRequest

	request.Unmarshal(&r)

	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

func ErrorSayHello(_ context.Context, _ *rpc.Request) (rpc.Response, error) {
	return nil, status.Error(http.StatusServiceUnavailable, "ohh no")
}
