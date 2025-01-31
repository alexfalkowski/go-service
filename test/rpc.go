package test

import (
	"cmp"
	"context"
	"errors"
	"net/http"

	"github.com/alexfalkowski/go-service/meta"
	nc "github.com/alexfalkowski/go-service/net/http/context"
	h "github.com/alexfalkowski/go-service/net/http/status"
	v1 "github.com/alexfalkowski/go-service/test/greet/v1"
	"google.golang.org/grpc/codes"
	g "google.golang.org/grpc/status"
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
func NoContent(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

// SuccessSayHello for test.
func SuccessSayHello(ctx context.Context, r *Request) (*Response, error) {
	req := nc.Request(ctx)
	_ = nc.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), r.Name)
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, ""), Greeting: s}, nil
}

// SuccessProtobufSayHello for test.
func SuccessProtobufSayHello(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

// ErrorSayHello for test.
func ErrorSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, h.Error(http.StatusInternalServerError, "ohh no")
}

// ErrorNotMappedSayHello for test.
//
//nolint:err113
func ErrorNotMappedSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, errors.New("ohh no")
}

// ErrorsProtobufSayHello for test.
func ErrorsProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, g.Error(codes.Internal, "ohh no")
}

// ErrorsNotMappedProtobufSayHello for test.
//
//nolint:err113
func ErrorsNotMappedProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, errors.New("ohh no")
}
