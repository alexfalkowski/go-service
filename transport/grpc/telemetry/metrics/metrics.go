package metrics

import (
	"context"
	"errors"
	"io"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/transport/strings"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcType string

const (
	unary        grpcType = "unary"
	clientStream grpcType = "client_stream"
	serverStream grpcType = "server_stream"
	bidiStream   grpcType = "bidi_stream"
)

func streamRPCType(info *grpc.StreamServerInfo) grpcType {
	if info.IsClientStream && !info.IsServerStream {
		return clientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return serverStream
	}

	return bidiStream
}

func clientStreamType(desc *grpc.StreamDesc) grpcType {
	if desc.ClientStreams && !desc.ServerStreams {
		return clientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return serverStream
	}

	return bidiStream
}

// NewServer for metrics.
func NewServer(meter metric.Meter) (*Server, error) {
	started, err := meter.Int64Counter("grpc_server_started_total", metric.WithDescription("Total number of RPCs started on the server."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Int64Counter("grpc_server_msg_received_total", metric.WithDescription("Total number of RPC messages received on the server."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Int64Counter("grpc_server_msg_sent_total", metric.WithDescription("Total number of RPC messages sent by the server."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("grpc_server_handled_total",
		metric.WithDescription("Total number of RPCs completed on the server, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("grpc_server_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of gRPC that had been application-level handled by the server."))
	if err != nil {
		return nil, err
	}

	s := &Server{
		started: started, received: received, sent: sent,
		handled: handled, handledHist: handledHist,
	}

	return s, nil
}

// Server for metrics.
type Server struct {
	started     metric.Int64Counter
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram
}

// UnaryInterceptor for metrics.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		opts := metric.WithAttributes(
			attribute.Key("grpc_type").String(string(unary)),
			attribute.Key("grpc_service").String(service),
			attribute.Key("grpc_method").String(method),
		)

		s.started.Add(ctx, 1, opts)
		s.received.Add(ctx, 1, opts)

		resp, err := handler(ctx, req)
		if err == nil {
			s.sent.Add(ctx, 1, opts)
		}

		s.handled.Add(ctx, 1, opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

		return resp, err
	}
}

// StreamInterceptor for metrics.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		opts := metric.WithAttributes(
			attribute.Key("grpc_type").String(string(streamRPCType(info))),
			attribute.Key("grpc_service").String(service),
			attribute.Key("grpc_method").String(method),
		)

		serverStream := &monitoredServerStream{
			opts: opts, received: s.received, sent: s.sent, handled: s.handled, handledHist: s.handledHist,
			ServerStream: stream,
		}

		err := handler(srv, serverStream)
		ctx := stream.Context()

		s.handled.Add(ctx, 1, opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

		return err
	}
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	opts        metric.MeasurementOption
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	grpc.ServerStream
}

func (s *monitoredServerStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
	s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

	return err
}

//nolint:dupl
func (s *monitoredServerStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	err := s.ServerStream.RecvMsg(m)
	if err != nil {
		if errors.Is(err, io.EOF) {
			s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(codes.OK.String())))
			s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

			return err
		}

		s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

		return err
	}

	s.received.Add(ctx, 1, s.opts)

	return err
}

// NewClient for metrics.
func NewClient(meter metric.Meter) (*Client, error) {
	started, err := meter.Int64Counter("grpc_client_started_total", metric.WithDescription("Total number of RPCs started on the client."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Int64Counter("grpc_client_msg_received_total", metric.WithDescription("Total number of RPC messages received on the client."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Int64Counter("grpc_client_msg_sent_total", metric.WithDescription("Total number of RPC messages sent by the client."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("grpc_client_handled_total",
		metric.WithDescription("Total number of RPCs completed on the client, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("grpc_client_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of gRPC that had been application-level handled by the client."))
	if err != nil {
		return nil, err
	}

	c := &Client{
		started: started, received: received, sent: sent,
		handled: handled, handledHist: handledHist,
	}

	return c, nil
}

// Client for metrics.
type Client struct {
	started     metric.Int64Counter
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram
}

// UnaryInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Unary RPCs.
func (c *Client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		o := metric.WithAttributes(
			attribute.Key("grpc_type").String(string(unary)),
			attribute.Key("grpc_service").String(service),
			attribute.Key("grpc_method").String(method),
		)

		c.started.Add(ctx, 1, o)
		c.sent.Add(ctx, 1, o)

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		if err == nil {
			c.received.Add(ctx, 1, o)
		}

		c.handled.Add(ctx, 1, o, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
		c.handledHist.Record(ctx, time.Since(start).Seconds(), o)

		return err
	}
}

// StreamInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Streaming RPCs.
func (c *Client) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		o := metric.WithAttributes(
			attribute.Key("grpc_type").String(string(clientStreamType(desc))),
			attribute.Key("grpc_service").String(service),
			attribute.Key("grpc_method").String(method),
		)

		clientStream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		if err != nil {
			c.handled.Add(ctx, 1, o, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
			c.handledHist.Record(ctx, time.Since(start).Seconds(), o)

			return nil, err
		}

		stream := &monitoredClientStream{
			opts:     o,
			received: c.received, sent: c.sent, handled: c.handled, handledHist: c.handledHist,
			ClientStream: clientStream,
		}

		return stream, nil
	}
}

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct {
	opts        metric.MeasurementOption
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	grpc.ClientStream
}

func (s *monitoredClientStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	err := s.ClientStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
	s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

	return err
}

//nolint:dupl
func (s *monitoredClientStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	err := s.ClientStream.RecvMsg(m)
	if err != nil {
		if errors.Is(err, io.EOF) {
			s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(codes.OK.String())))
			s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

			return err
		}

		s.handled.Add(ctx, 1, s.opts, metric.WithAttributes(attribute.Key("grpc_code").String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), s.opts)

		return err
	}

	s.received.Add(ctx, 1, s.opts)

	return nil
}
