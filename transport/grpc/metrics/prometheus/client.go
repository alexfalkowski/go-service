package prometheus

import (
	"errors"
	"io"
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

// ClientMetrics for prometheus.
type ClientMetrics struct {
	clientStartedCounter   *prometheus.CounterVec
	clientHandledCounter   *prometheus.CounterVec
	clientMsgReceived      *prometheus.CounterVec
	clientMsgSent          *prometheus.CounterVec
	clientHandledHistogram *prometheus.HistogramVec
	clientRecvHistogram    *prometheus.HistogramVec
	clientSendHistogram    *prometheus.HistogramVec
}

// NewClientMetrics for prometheus.
func NewClientMetrics(lc fx.Lifecycle, version version.Version) *ClientMetrics {
	labels := prometheus.Labels{"name": os.ExecutableName(), "version": string(version)}

	metrics := &ClientMetrics{
		clientStartedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_client_started_total",
				Help:        "Total number of RPCs started on the client.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		clientHandledCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_client_handled_total",
				Help:        "Total number of RPCs completed by the client, regardless of success or failure.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method", "grpc_code"}),
		clientMsgReceived: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_client_msg_received_total",
				Help:        "Total number of RPC messages received by the client.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		clientMsgSent: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "grpc_client_msg_sent_total",
				Help:        "Total number of gRPC messages sent by the client.",
				ConstLabels: labels,
			}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		clientHandledHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "grpc_client_handling_seconds",
			Help:        "Histogram of response latency (seconds) of the gRPC until it is finished by the application.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: labels,
		}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		clientRecvHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "grpc_client_msg_recv_handling_seconds",
			Help:        "Histogram of response latency (seconds) of the gRPC single message receive.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: labels,
		}, []string{"grpc_type", "grpc_service", "grpc_method"}),
		clientSendHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "grpc_client_msg_send_handling_seconds",
			Help:        "Histogram of response latency (seconds) of the gRPC single message send.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: labels,
		}, []string{"grpc_type", "grpc_service", "grpc_method"}),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prometheus.Register(metrics)
		},
		OnStop: func(ctx context.Context) error {
			prometheus.Unregister(metrics)

			return nil
		},
	})

	return metrics
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	m.clientStartedCounter.Describe(ch)
	m.clientHandledCounter.Describe(ch)
	m.clientMsgReceived.Describe(ch)
	m.clientMsgSent.Describe(ch)
	m.clientHandledHistogram.Describe(ch)
	m.clientRecvHistogram.Describe(ch)
	m.clientSendHistogram.Describe(ch)
}

// Collect is called by the prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ClientMetrics) Collect(ch chan<- prometheus.Metric) {
	m.clientStartedCounter.Collect(ch)
	m.clientHandledCounter.Collect(ch)
	m.clientMsgReceived.Collect(ch)
	m.clientMsgSent.Collect(ch)
	m.clientHandledHistogram.Collect(ch)
	m.clientRecvHistogram.Collect(ch)
	m.clientSendHistogram.Collect(ch)
}

// UnaryClientInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Unary RPCs.
func (m *ClientMetrics) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, fullMethod string, req, resp any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return invoker(ctx, fullMethod, req, resp, cc, opts...)
		}

		method := path.Base(fullMethod)
		monitor := newClientReporter(m, unary, service, method)
		monitor.SentMessage()

		err := invoker(ctx, fullMethod, req, resp, cc, opts...)
		if err == nil {
			monitor.ReceivedMessage()
		}

		st, _ := status.FromError(err)
		monitor.Handled(st.Code())

		return err
	}
}

// StreamClientInterceptor is a gRPC client-side interceptor that provides prometheus monitoring for Streaming RPCs.
func (m *ClientMetrics) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		service := path.Dir(fullMethod)[1:]
		if strings.IsHealth(service) {
			return streamer(ctx, desc, cc, fullMethod, opts...)
		}

		method := path.Base(fullMethod)
		monitor := newClientReporter(m, clientStreamType(desc), service, method)

		clientStream, err := streamer(ctx, desc, cc, fullMethod, opts...)
		if err != nil {
			st, _ := status.FromError(err)
			monitor.Handled(st.Code())

			return nil, err
		}

		return &monitoredClientStream{clientStream, monitor}, nil
	}
}

func clientStreamType(desc *grpc.StreamDesc) grpcType {
	if desc.ClientStreams && !desc.ServerStreams {
		return clientStream
	} else if !desc.ClientStreams && desc.ServerStreams {
		return serverStream
	}

	return bidiStream
}

// monitoredClientStream wraps grpc.ClientStream allowing each Sent/Recv of message to increment counters.
type monitoredClientStream struct {
	grpc.ClientStream
	monitor *clientReporter
}

func (s *monitoredClientStream) SendMsg(m any) error {
	timer := s.monitor.SendMessageTimer()
	err := s.ClientStream.SendMsg(m)

	timer.ObserveDuration()

	if err == nil {
		s.monitor.SentMessage()
	}

	return err
}

func (s *monitoredClientStream) RecvMsg(m any) error {
	timer := s.monitor.ReceiveMessageTimer()
	err := s.ClientStream.RecvMsg(m)

	timer.ObserveDuration()

	if err != nil {
		if errors.Is(err, io.EOF) {
			s.monitor.Handled(codes.OK)

			return err
		}

		st, _ := status.FromError(err)
		s.monitor.Handled(st.Code())

		return err
	}

	s.monitor.ReceivedMessage()

	return nil
}

type clientReporter struct {
	metrics     *ClientMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
}

func newClientReporter(m *ClientMetrics, rpcType grpcType, service, method string) *clientReporter {
	r := &clientReporter{metrics: m, rpcType: rpcType, startTime: time.Now(), serviceName: service, methodName: method}
	r.metrics.clientStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()

	return r
}

// timer is a helper interface to time functions.
type timer interface {
	ObserveDuration() time.Duration
}

func (r *clientReporter) ReceiveMessageTimer() timer {
	hist := r.metrics.clientRecvHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName)

	return prometheus.NewTimer(hist)
}

func (r *clientReporter) ReceivedMessage() {
	r.metrics.clientMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) SendMessageTimer() timer {
	hist := r.metrics.clientSendHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName)

	return prometheus.NewTimer(hist)
}

func (r *clientReporter) SentMessage() {
	r.metrics.clientMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Inc()
}

func (r *clientReporter) Handled(code codes.Code) {
	r.metrics.clientHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	r.metrics.clientHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName).Observe(time.Since(r.startTime).Seconds())
}
