package test

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
)

// MetaServerStream is a grpc.ServerStream test double that records response headers.
type MetaServerStream struct {
	grpc.ServerStream
	Header meta.Map
	Ctx    context.Context
}

// Context returns Ctx.
func (s *MetaServerStream) Context() context.Context {
	return s.Ctx
}

// SetHeader records md after appending a service-version value.
func (s *MetaServerStream) SetHeader(md meta.Map) error {
	md.Append("service-version", "v2")
	s.Header = md

	return nil
}

// NewWatchStream returns a WatchStream with a buffered response channel.
func NewWatchStream(ctx context.Context) *WatchStream {
	return &WatchStream{Ctx: ctx, Responses: make(chan *health.Response, 4)}
}

// WatchStream is a grpc.ServerStream test double for health watch tests.
type WatchStream struct {
	grpc.ServerStream
	Ctx       context.Context
	Responses chan *health.Response
}

// Context returns Ctx.
func (w *WatchStream) Context() context.Context {
	return w.Ctx
}

// Send records resp or returns the context error if canceled.
func (w *WatchStream) Send(resp *health.Response) error {
	select {
	case <-w.Ctx.Done():
		return w.Ctx.Err()
	case w.Responses <- resp:
		return nil
	}
}

// SetHeader implements grpc.ServerStream.
func (*WatchStream) SetHeader(meta.Map) error {
	return nil
}

// SendHeader implements grpc.ServerStream.
func (*WatchStream) SendHeader(meta.Map) error {
	return nil
}

// SetTrailer implements grpc.ServerStream.
func (*WatchStream) SetTrailer(meta.Map) {}

// SendMsg implements grpc.ServerStream.
func (*WatchStream) SendMsg(any) error {
	return nil
}

// RecvMsg implements grpc.ServerStream.
func (*WatchStream) RecvMsg(any) error {
	return nil
}
