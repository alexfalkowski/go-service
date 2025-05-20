package test

import (
	"cmp"
	"context"
	"net/http"

	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/meta"
	hm "github.com/alexfalkowski/go-service/v2/net/http/meta"
	h "github.com/alexfalkowski/go-service/v2/net/http/status"
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
	req := hm.Request(ctx)
	_ = hm.Response(ctx)
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
	return nil, h.Error(http.StatusInternalServerError, ErrFailed.Error())
}

// ErrorNotMappedSayHello for test.
func ErrorNotMappedSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, ErrFailed
}

// ErrorsProtobufSayHello for test.
func ErrorsProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, g.Error(codes.Internal, ErrFailed.Error())
}

// ErrorsNotMappedProtobufSayHello for test.
func ErrorsNotMappedProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, ErrFailed
}

// ErrorsInternalProtobufSayHello for test.
func ErrorsInternalProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, ErrInternal
}
