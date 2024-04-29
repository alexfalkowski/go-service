package metrics

import (
	"context"
	"errors"
	"io"
	"path"
	"time"

	"github.com/alexfalkowski/go-service/telemetry/metrics"
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
		handled: handled, handledHist: handledHist,
	}

	return s
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
			kindAttribute.String(string(unary)),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		s.started.Add(ctx, 1, opts)
		s.received.Add(ctx, 1, opts)

		resp, err := handler(ctx, req)
		if err == nil {
			s.sent.Add(ctx, 1, opts)
		}

		s.handled.Add(ctx, 1, opts, metric.WithAttributes(codeAttribute.String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

		return resp, err
	}
}

// StreamInterceptor for metrics.
func (s *Server) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, st grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, st)
		}

		start := time.Now()
		method := path.Base(info.FullMethod)
		opts := metric.WithAttributes(
			kindAttribute.String(string(stream)),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		stream := &serverStream{
			opts: opts, received: s.received, sent: s.sent, handled: s.handled, handledHist: s.handledHist,
			ServerStream: st,
		}

		err := handler(srv, stream)
		ctx := st.Context()

		s.handled.Add(ctx, 1, opts, metric.WithAttributes(codeAttribute.String(status.Code(err).String())))
		s.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

		return err
	}
}

// serverStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type serverStream struct {
	opts        metric.MeasurementOption
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	grpc.ServerStream
}

func (s *serverStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handleStream(ctx, s.handled, s.handledHist, s.opts, status.Code(err), start)

	return err
}

func (s *serverStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ServerStream.Context()

	if err := s.ServerStream.RecvMsg(m); err != nil {
		if errors.Is(err, io.EOF) {
			handleStream(ctx, s.handled, s.handledHist, s.opts, codes.OK, start)

			return err
		}

		handleStream(ctx, s.handled, s.handledHist, s.opts, status.Code(err), start)

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
		handled: handled, handledHist: handledHist,
	}

	return c
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
			kindAttribute.String(string(unary)),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		c.started.Add(ctx, 1, o)
		c.sent.Add(ctx, 1, o)

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		if err == nil {
			c.received.Add(ctx, 1, o)
		}

		c.handled.Add(ctx, 1, o, metric.WithAttributes(codeAttribute.String(status.Code(err).String())))
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
			kindAttribute.String(string(stream)),
			serviceAttribute.String(service),
			methodAttribute.String(method),
		)

		stream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		if err != nil {
			c.handled.Add(ctx, 1, o, metric.WithAttributes(codeAttribute.String(status.Code(err).String())))
			c.handledHist.Record(ctx, time.Since(start).Seconds(), o)

			return nil, err
		}

		st := &clientStream{
			opts:     o,
			received: c.received, sent: c.sent, handled: c.handled, handledHist: c.handledHist,
			ClientStream: stream,
		}

		return st, nil
	}
}

// clientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type clientStream struct {
	opts        metric.MeasurementOption
	received    metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	grpc.ClientStream
}

func (s *clientStream) SendMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	err := s.ClientStream.SendMsg(m)
	if err == nil {
		s.sent.Add(ctx, 1, s.opts)
	}

	handleStream(ctx, s.handled, s.handledHist, s.opts, status.Code(err), start)

	return err
}

func (s *clientStream) RecvMsg(m any) error {
	start := time.Now()
	ctx := s.ClientStream.Context()

	if err := s.ClientStream.RecvMsg(m); err != nil {
		if errors.Is(err, io.EOF) {
			handleStream(ctx, s.handled, s.handledHist, s.opts, codes.OK, start)

			return err
		}

		handleStream(ctx, s.handled, s.handledHist, s.opts, status.Code(err), start)

		return err
	}

	s.received.Add(ctx, 1, s.opts)

	return nil
}

func handleStream(ctx context.Context, h metric.Int64Counter, hs metric.Float64Histogram, o metric.MeasurementOption, c codes.Code, s time.Time) {
	h.Add(ctx, 1, o, metric.WithAttributes(codeAttribute.String(c.String())))
	hs.Record(ctx, time.Since(s).Seconds(), o)
}
