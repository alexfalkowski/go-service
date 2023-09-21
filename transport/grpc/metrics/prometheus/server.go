package prometheus

import (
	"path"
	"time"

	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/transport/strings"
	"github.com/alexfalkowski/go-service/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerCollector for prometheus.
type ServerCollector struct {
	started          *prometheus.CounterVec
	handled          *prometheus.CounterVec
	received         *prometheus.CounterVec
	sent             *prometheus.CounterVec
	handledHistogram *prometheus.HistogramVec
}

// NewServerCollector for prometheus.
func NewServerCollector(lc fx.Lifecycle, version version.Version) *ServerCollector {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	collector := &ServerCollector{
		started: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_server_started_total",
				Help:        "Total number of RPCs started on the server.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		handled: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_server_handled_total",
				Help:        "Total number of RPCs completed on the server, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}),
		received: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_server_msg_received_total",
				Help:        "Total number of RPC messages received on the server.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		sent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_server_msg_sent_total",
				Help:        "Total number of gRPC messages sent by the server.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		handledHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "grpc_server_handling_seconds",
				Help:        "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"},
		),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(collector)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(collector)

			return nil
		},
	})

	return collector
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerCollector) Describe(ch chan<- *prometheus.Desc) {
	m.started.Describe(ch)
	m.handled.Describe(ch)
	m.received.Describe(ch)
	m.sent.Describe(ch)
	m.handledHistogram.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerCollector) Collect(ch chan<- prometheus.Metric) {
	m.started.Collect(ch)
	m.handled.Collect(ch)
	m.received.Collect(ch)
	m.sent.Collect(ch)
	m.handledHistogram.Collect(ch)
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ServerCollector) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(ctx, req)
		}

		method := path.Base(info.FullMethod)
		monitor := newServerReporter(m, unary, service, method)
		monitor.ReceivedMessage()

		resp, err := handler(ctx, req)
		st, _ := status.FromError(err)
		monitor.Handled(st.Code())

		if err == nil {
			monitor.SentMessage()
		}

		return resp, err
	}
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func (m *ServerCollector) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		service := path.Dir(info.FullMethod)[1:]
		if strings.IsHealth(service) {
			return handler(srv, stream)
		}

		method := path.Base(info.FullMethod)
		monitor := newServerReporter(m, streamRPCType(info), service, method)
		err := handler(srv, &monitoredServerStream{stream, monitor})
		st, _ := status.FromError(err)

		monitor.Handled(st.Code())

		return err
	}
}

func streamRPCType(info *grpc.StreamServerInfo) grpcType {
	if info.IsClientStream && !info.IsServerStream {
		return clientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return serverStream
	}

	return bidiStream
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	grpc.ServerStream
	monitor *serverReporter
}

func (s *monitoredServerStream) SendMsg(m any) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.monitor.SentMessage()
	}

	return err
}

func (s *monitoredServerStream) RecvMsg(m any) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.monitor.ReceivedMessage()
	}

	return err
}

type serverReporter struct {
	metrics     *ServerCollector
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
}

func newServerReporter(m *ServerCollector, rpcType grpcType, service, method string) *serverReporter {
	r := &serverReporter{metrics: m, rpcType: rpcType, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.started.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()

	return r
}

func (r *serverReporter) ReceivedMessage() {
	r.metrics.received.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) SentMessage() {
	r.metrics.sent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *serverReporter) Handled(code codes.Code) {
	r.metrics.handled.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	r.metrics.handledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}
