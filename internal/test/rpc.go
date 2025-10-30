package test

import (
	"cmp"

	"github.com/alexfalkowski/go-service/v2/context"
	v1 "github.com/alexfalkowski/go-service/v2/internal/test/greet/v1"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	g "github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/meta"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	h "github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/strings"
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
	return &Response{}, nil
}

// SuccessSayHello for test.
func SuccessSayHello(ctx context.Context, r *Request) (*Response, error) {
	req := meta.Request(ctx)
	_ = meta.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), r.Name)
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, strings.Empty), Greeting: s}, nil
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

func (w *World) registerRPC() {
	rpc.Register(rpc.RegisterParams{
		Mux:     w.ServeMux,
		Pool:    Pool,
		Content: Content,
	})
}
