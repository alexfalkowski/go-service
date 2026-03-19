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
)

// Request is the simple RPC request payload used by HTTP RPC and REST helper handlers.
type Request struct {
	Name string
}

// Response is the simple RPC response payload used by HTTP RPC and REST helper handlers.
type Response struct {
	Meta     meta.Map
	Greeting string
}

// NoContent returns an empty response value and no error.
func NoContent(_ context.Context, _ *Request) (*Response, error) {
	return &Response{}, nil
}

// SuccessSayHello resolves the greeting name from the query string first, then the request body, and returns request metadata.
func SuccessSayHello(ctx context.Context, r *Request) (*Response, error) {
	req := meta.Request(ctx)
	_ = meta.Response(ctx)
	name := cmp.Or(req.URL.Query().Get("name"), r.Name)
	s := "Hello " + name

	return &Response{Meta: meta.CamelStrings(ctx, meta.NoPrefix), Greeting: s}, nil
}

// SuccessProtobufSayHello returns a protobuf greeting response for the supplied request name.
func SuccessProtobufSayHello(_ context.Context, r *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return &v1.SayHelloResponse{Message: "Hello " + r.GetName()}, nil
}

// ErrorSayHello returns a mapped HTTP status error with an internal server status code.
func ErrorSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, h.Error(http.StatusInternalServerError, ErrFailed.Error())
}

// ErrorNotMappedSayHello returns ErrFailed directly so callers can exercise fallback error mapping.
func ErrorNotMappedSayHello(_ context.Context, _ *Request) (*Response, error) {
	return nil, ErrFailed
}

// ErrorsProtobufSayHello returns a mapped gRPC internal status error.
func ErrorsProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, g.Error(codes.Internal, ErrFailed.Error())
}

// ErrorsNotMappedProtobufSayHello returns ErrFailed directly so callers can exercise fallback gRPC error mapping.
func ErrorsNotMappedProtobufSayHello(_ context.Context, _ *v1.SayHelloRequest) (*v1.SayHelloResponse, error) {
	return nil, ErrFailed
}

// ErrorsInternalProtobufSayHello returns ErrInternal, which already implements the HTTP status coder contract.
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
