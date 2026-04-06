package health_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
	v1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

func TestCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestInvalidCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestNotFoundCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: "bob"}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	_, err = client.Check(ctx, req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))

	world.RequireStop()
}

func TestIgnoreAuthCheck(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(t.Context(), req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	world.RequireStop()
}

func TestList(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthListRequest{}

	resp, err := client.List(ctx, req)
	require.NoError(t, err)

	expected := map[string]*v1.HealthCheckResponse{
		test.Name.String(): {Status: v1.HealthCheckResponse_SERVING},
	}
	require.Equal(t, expected, resp.GetStatuses())

	world.RequireStop()
}

func TestWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)

	world.RequireStop()
}

func TestInvalidWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)

	world.RequireStop()
}

func TestNotFoundWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("500"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: "bob"}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)
	require.Equal(t, v1.HealthCheckResponse_SERVICE_UNKNOWN, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)

	world.RequireStop()
}

func TestIgnoreAuthWatch(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"), test.WithWorldToken(nil, test.NewVerifier("test")), test.WithWorldGRPC())
	world.Register()

	so := world.HealthServer(test.Name.String(), test.StatusURL("200"))

	err := so.Observe(test.Name.String(), "grpc", "http")
	require.NoError(t, err)

	server := health.NewServer(health.ServerParams{Server: so})
	health.Register(health.RegisterParams{Registrar: world.GRPCServer.ServiceRegistrar(), Server: server})

	world.RequireStart()
	time.Sleep(1 * time.Second)

	conn := world.NewGRPC()
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)

	world.RequireStop()
}

func TestWatchStatusChanges(t *testing.T) {
	world := test.NewWorld(t, test.WithWorldTelemetry("otlp"))
	world.Register()

	var unhealthy sync.Bool
	probe := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if unhealthy.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer probe.Close()

	so := world.HealthServer(test.Name.String(), probe.URL)
	require.NoError(t, so.Observe(test.Name.String(), "grpc", "http"))

	world.RequireStart()
	defer world.RequireStop()

	watcher := health.NewServer(health.ServerParams{Server: so})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	stream := newWatchStream(ctx)
	errCh := make(chan error, 1)
	go func() {
		errCh <- watcher.Watch(&v1.HealthCheckRequest{Service: test.Name.String()}, stream)
	}()

	resp := requireWatchResponse(t, stream.responses)
	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())

	unhealthy.Store(true)

	require.Eventually(t, func() bool {
		select {
		case resp = <-stream.responses:
			return resp.GetStatus() == v1.HealthCheckResponse_NOT_SERVING
		default:
			return false
		}
	}, time.Second, 25*time.Millisecond)

	cancel()

	select {
	case err := <-errCh:
		require.Error(t, err)
		require.Equal(t, codes.Canceled, status.Code(err))
	case <-time.After(time.Second):
		require.FailNow(t, "watch stream did not stop after cancellation")
	}
}

func requireWatchStaysOpen(t *testing.T, cancel context.CancelFunc, wc v1.Health_WatchClient) {
	t.Helper()

	errCh := make(chan error, 1)
	go func() {
		_, err := wc.Recv()
		errCh <- err
	}()

	select {
	case err := <-errCh:
		require.FailNow(t, "watch stream closed unexpectedly", err.Error())
	case <-time.After(150 * time.Millisecond):
	}

	cancel()

	select {
	case err := <-errCh:
		require.Error(t, err)
		require.Equal(t, codes.Canceled, status.Code(err))
	case <-time.After(time.Second):
		require.FailNow(t, "watch stream did not stop after cancellation")
	}
}

func requireWatchResponse(t *testing.T, responses <-chan *v1.HealthCheckResponse) *v1.HealthCheckResponse {
	t.Helper()

	select {
	case resp := <-responses:
		return resp
	case <-time.After(time.Second):
		require.FailNow(t, "timed out waiting for watch response")
		return nil
	}
}

func newWatchStream(ctx context.Context) *watchStream {
	return &watchStream{ctx: ctx, responses: make(chan *v1.HealthCheckResponse, 4)}
}

type watchStream struct {
	grpc.ServerStream
	ctx       context.Context
	responses chan *v1.HealthCheckResponse
}

func (w *watchStream) Context() context.Context {
	return w.ctx
}

func (w *watchStream) Send(resp *v1.HealthCheckResponse) error {
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	case w.responses <- resp:
		return nil
	}
}

func (*watchStream) SetHeader(metadata.MD) error {
	return nil
}

func (*watchStream) SendHeader(metadata.MD) error {
	return nil
}

func (*watchStream) SetTrailer(metadata.MD) {}

func (*watchStream) SendMsg(any) error {
	return nil
}

func (*watchStream) RecvMsg(any) error {
	return nil
}
