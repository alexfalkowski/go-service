package metrics

import (
	"context"
	"errors"
	"io"
	"path"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/strings"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewServer for metrics.
func NewServer(meter metric.Meter) *Server {
	started := metrics.MustInt64Counter(meter, "grpc_server_started_total", "Total number of RPCs started on the server.")
	received := metrics.MustInt64Counter(meter, "grpc_server_msg_received_total", "Total number of RPC messages received on the server.")
	sent := metrics.MustInt64Counter(meter, "grpc_server_msg_sent_total", "Total number of RPC messages sent by the server.")
	handled := metrics.MustInt64Counter(meter, "grpc_server_handled_total", "Total number of RPCs completed on the server, regardless of success or failure.")
	handledHist := metrics.MustFloat64Histogram(meter, "grpc_server_handling_seconds",
		"Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.")

	s := &Server{
		started: started, received: received, sent: sent,
		handled: handled, handledHistogram: handledHist,
	}

	return s
}

// Server for metrics.
type Server struct {
	started          metric.Int64Counter
	received         metric.Int64Counter
	sent             metric.Int64Counter
	handled          metric.Int64Counter
	handledHistogram metric.Float64Histogram
}

// UnaryInterceptor for metrics.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(ctx, req)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		opts := metric.WithAttributes(
			kindAttribute.String(UnaryKind.String()),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		s.started.Add(ctx, 1, opts)
		s.received.Add(ctx, 1, opts)

		resp, err := handler(ctx, req)
		if err == nil {
			s.sent.Add(ctx, 1, opts)
		}

		handle(ctx, s.handled, s.handledHistogram, opts, status.Code(err), start)

		return resp, err
	}
}

// StreamInterceptor for metrics.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsObservable(service) {
			return handler(srv, stream)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		opts := metric.WithAttributes(
			kindAttribute.String(StreamKind.String()),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)
		err := handler(srv, s.Stream(stream, opts))
		ctx := stream.Context()

		handle(ctx, s.handled, s.handledHistogram, opts, status.Code(err), start)

		return err
	}
}

// Stream for the server.
func (s *Server) Stream(stream grpc.ServerStream, opts metric.MeasurementOption) grpc.ServerStream {
	return &ServerStream{
		opts:             opts,
		received:         s.received,
		sent:             s.sent,
		handled:          s.handled,
		handledHistogram: s.handledHistogram,
		ServerStream:     stream,
	}
}

// ServerStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type ServerStream struct {
	opts             metric.MeasurementOption
	received         metric.Int64Counter
	sent             metric.Int64Counter
	handled          metric.Int64Counter
	handledHistogram metric.Float64Histogram

	grpc.ServerStream
}

func (s *ServerStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

	return err
}

func (s *ServerStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	if err := s.ServerStream.RecvMsg(m); err != nil {
		if errors.Is(err, io.EOF) {
			handle(ctx, s.handled, s.handledHistogram, s.opts, codes.OK, start)

			return err
		}

		handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

		return err
	}

	s.received.Add(ctx, 1, s.opts)

	return nil
}

// NewClient for metrics.
func NewClient(meter metric.Meter) *Client {
	started := metrics.MustInt64Counter(meter, "grpc_client_started_total", "Total number of RPCs started on the client.")
	received := metrics.MustInt64Counter(meter, "grpc_client_msg_received_total", "Total number of RPC messages received on the client.")
	sent := metrics.MustInt64Counter(meter, "grpc_client_msg_sent_total", "Total number of RPC messages sent by the client.")
	handled := metrics.MustInt64Counter(meter, "grpc_client_handled_total", "Total number of RPCs completed on the client, regardless of success or failure.")
	handledHist := metrics.MustFloat64Histogram(meter, "grpc_client_handling_seconds",
		"Histogram of response latency (seconds) of gRPC that had been application-level handled by the client.")

	c := &Client{
		started: started, received: received, sent: sent,
		handled: handled, handledHistogram: handledHist,
	}

	return c
}

// Client for metrics.
type Client struct {
	started          metric.Int64Counter
	received         metric.Int64Counter
	sent             metric.Int64Counter
	handled          metric.Int64Counter
	handledHistogram metric.Float64Histogram
}

// UnaryInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Unary RPCs.
func (c *Client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		measurement := metric.WithAttributes(
			kindAttribute.String(UnaryKind.String()),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		c.started.Add(ctx, 1, measurement)
		c.sent.Add(ctx, 1, measurement)

		err := invoker(ctx, fullMethod, req, resp, conn, opts...)
		if err == nil {
			c.received.Add(ctx, 1, measurement)
		}

		handle(ctx, c.handled, c.handledHistogram, measurement, status.Code(err), start)

		return err
	}
}

// StreamInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Streaming RPCs.
func (c *Client) StreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsObservable(service) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		start := time.Now()
		method := path.Base(fullMethod)
		measurement := metric.WithAttributes(
			kindAttribute.String(StreamKind.String()),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		stream, err := streamer(ctx, desc, conn, fullMethod, opts...)
		if err != nil {
			handle(ctx, c.handled, c.handledHistogram, measurement, status.Code(err), start)

			return nil, err
		}

		return c.Stream(stream, measurement), nil
	}
}

// Stream fpr client.
func (c *Client) Stream(stream grpc.ClientStream, opts metric.MeasurementOption) grpc.ClientStream {
	return &ClientStream{
		opts:             opts,
		received:         c.received,
		sent:             c.sent,
		handled:          c.handled,
		handledHistogram: c.handledHistogram,
		ClientStream:     stream,
	}
}

// ClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type ClientStream struct {
	opts             metric.MeasurementOption
	received         metric.Int64Counter
	sent             metric.Int64Counter
	handled          metric.Int64Counter
	handledHistogram metric.Float64Histogram

	grpc.ClientStream
}

func (s *ClientStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	err := s.ClientStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

	return err
}

func (s *ClientStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	if err := s.ClientStream.RecvMsg(m); err != nil {
		if errors.Is(err, io.EOF) {
			handle(ctx, s.handled, s.handledHistogram, s.opts, codes.OK, start)

			return err
		}

		handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

		return err
	}

	s.received.Add(ctx, 1, s.opts)

	return nil
}

func handle(ctx context.Context, h metric.Int64Counter, hs metric.Float64Histogram, o metric.MeasurementOption, c codes.Code, s time.Time) {
	h.Add(ctx, 1, o, metric.WithAttributes(codeAttribute.String(c.String())))
	hs.Record(ctx, time.Since(s).Seconds(), o)
}
