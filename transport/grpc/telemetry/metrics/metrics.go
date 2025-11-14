package metrics

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/strings"
)

const (
	kindAttribute    = attributes.Key("service_kind")
	nameAttribute    = attributes.Key("service_name")
	serviceAttribute = attributes.Key("grpc_service")
	methodAttribute  = attributes.Key("grpc_method")
	codeAttribute    = attributes.Key("grpc_code")
)

// Meter is an alias for metrics.Meter.
type Meter = metrics.Meter

// NewServer for metrics.
func NewServer(name env.Name, meter *Meter) *Server {
	started := meter.MustInt64Counter("grpc_server_started_total", "Total number of RPCs started on the server.")
	received := meter.MustInt64Counter("grpc_server_msg_received_total", "Total number of RPC messages received on the server.")
	sent := meter.MustInt64Counter("grpc_server_msg_sent_total", "Total number of RPC messages sent by the server.")
	handled := meter.MustInt64Counter("grpc_server_handled_total", "Total number of RPCs completed on the server, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("grpc_server_handling_seconds",
		"Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.")

	return &Server{
		name:             name,
		started:          started,
		received:         received,
		sent:             sent,
		handled:          handled,
		handledHistogram: handledHist,
	}
}

// Server for metrics.
type Server struct {
	started          metrics.Int64Counter
	received         metrics.Int64Counter
	sent             metrics.Int64Counter
	handled          metrics.Int64Counter
	handledHistogram metrics.Float64Histogram
	name             env.Name
}

// UnaryInterceptor for metrics.
func (s *Server) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.IsObservable(info.FullMethod) {
			return handler(ctx, req)
		}

		service, method := strings.SplitServiceMethod(info.FullMethod)
		start := time.Now()
		opts := metrics.WithAttributes(
			kindAttribute.String(UnaryKind.String()),
			nameAttribute.String(s.name.String()),
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
		if strings.IsObservable(info.FullMethod) {
			return handler(srv, stream)
		}

		service, method := strings.SplitServiceMethod(info.FullMethod)
		start := time.Now()
		opts := metrics.WithAttributes(
			kindAttribute.String(StreamKind.String()),
			nameAttribute.String(s.name.String()),
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
func (s *Server) Stream(stream grpc.ServerStream, opts metrics.MeasurementOption) grpc.ServerStream {
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
	opts             metrics.MeasurementOption
	received         metrics.Int64Counter
	sent             metrics.Int64Counter
	handled          metrics.Int64Counter
	handledHistogram metrics.Float64Histogram
	grpc.ServerStream
}

func (s *ServerStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.Context()

	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

	return err
}

func (s *ServerStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.Context()

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
func NewClient(name env.Name, meter *Meter) *Client {
	started := meter.MustInt64Counter("grpc_client_started_total", "Total number of RPCs started on the client.")
	received := meter.MustInt64Counter("grpc_client_msg_received_total", "Total number of RPC messages received on the client.")
	sent := meter.MustInt64Counter("grpc_client_msg_sent_total", "Total number of RPC messages sent by the client.")
	handled := meter.MustInt64Counter("grpc_client_handled_total", "Total number of RPCs completed on the client, regardless of success or failure.")
	handledHist := meter.MustFloat64Histogram("grpc_client_handling_seconds",
		"Histogram of response latency (seconds) of gRPC that had been application-level handled by the client.")

	return &Client{
		name:             name,
		started:          started,
		received:         received,
		sent:             sent,
		handled:          handled,
		handledHistogram: handledHist,
	}
}

// Client for metrics.
type Client struct {
	started          metrics.Int64Counter
	received         metrics.Int64Counter
	sent             metrics.Int64Counter
	handled          metrics.Int64Counter
	handledHistogram metrics.Float64Histogram
	name             env.Name
}

// UnaryInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Unary RPCs.
func (c *Client) UnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if strings.IsObservable(fullMethod) {
			return invoker(ctx, fullMethod, req, resp, conn, opts...)
		}

		service, method := strings.SplitServiceMethod(fullMethod)
		start := time.Now()
		measurement := metrics.WithAttributes(
			kindAttribute.String(UnaryKind.String()),
			nameAttribute.String(c.name.String()),
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
		if strings.IsObservable(fullMethod) {
			return streamer(ctx, desc, conn, fullMethod, opts...)
		}

		service, method := strings.SplitServiceMethod(fullMethod)
		start := time.Now()
		measurement := metrics.WithAttributes(
			kindAttribute.String(StreamKind.String()),
			nameAttribute.String(c.name.String()),
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
func (c *Client) Stream(stream grpc.ClientStream, opts metrics.MeasurementOption) grpc.ClientStream {
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
	opts             metrics.MeasurementOption
	received         metrics.Int64Counter
	sent             metrics.Int64Counter
	handled          metrics.Int64Counter
	handledHistogram metrics.Float64Histogram
	grpc.ClientStream
}

func (s *ClientStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.Context()

	err := s.ClientStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handle(ctx, s.handled, s.handledHistogram, s.opts, status.Code(err), start)

	return err
}

func (s *ClientStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.Context()

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

func handle(ctx context.Context, h metrics.Int64Counter, hs metrics.Float64Histogram, o metrics.MeasurementOption, c codes.Code, s time.Time) {
	h.Add(ctx, 1, o, metrics.WithAttributes(codeAttribute.String(c.String())))
	hs.Record(ctx, time.Since(s).Seconds(), o)
}
