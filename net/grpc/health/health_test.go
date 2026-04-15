package health_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-health/v2/server"
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
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
	)
	requireGRPCReady(t, world)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestInvalidCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)
	requireObservedHealth(t, world.GRPCHealth, test.Name.String(), false)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_NOT_SERVING, resp.GetStatus())
}

func TestNotFoundCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: "bob"}

	md := metadata.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := metadata.NewOutgoingContext(t.Context(), md)

	_, err := client.Check(ctx, req)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestIgnoreAuthCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(nil, test.NewVerifier("test")),
	)
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthCheckRequest{Service: test.Name.String()}

	resp, err := client.Check(t.Context(), req)
	require.NoError(t, err)

	require.Equal(t, v1.HealthCheckResponse_SERVING, resp.GetStatus())
}

func TestList(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
	)
	requireGRPCReady(t, world)

	ctx := t.Context()
	ctx = meta.WithRequestID(ctx, meta.String("test-id"))
	ctx = meta.WithUserAgent(ctx, meta.String("test-user-agent"))

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := v1.NewHealthClient(conn)
	req := &v1.HealthListRequest{}

	resp, err := client.List(ctx, req)
	require.NoError(t, err)

	expected := map[string]*v1.HealthCheckResponse{
		test.Name.String(): {Status: v1.HealthCheckResponse_SERVING},
	}
	require.Equal(t, expected, resp.GetStatuses())
}

func TestWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 10)),
	)
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
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
}

func TestInvalidWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)
	requireObservedHealth(t, world.GRPCHealth, test.Name.String(), false)

	conn := requireGRPCConn(t, world)
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
}

func TestNotFoundWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
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
}

func TestIgnoreAuthWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldToken(nil, test.NewVerifier("test")),
	)
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
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
}

func TestWatchStatusChanges(t *testing.T) {
	var unhealthy sync.Bool
	probe := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if unhealthy.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer probe.Close()

	world := newGRPCHealthWorld(t, probe.URL, test.WithWorldTelemetry("otlp"))

	watcher := health.NewServer(health.ServerParams{Server: world.GRPCHealth})
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
	}, time.Second.Duration(), (25 * time.Millisecond).Duration())

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

func requireGRPCReady(t *testing.T, world *test.World) {
	t.Helper()

	conn, err := test.Connect(t.Context(), world.TransportConfig.GRPC.Address)
	require.NoError(t, err)
	require.NoError(t, conn.Close())
}

func requireGRPCConn(t *testing.T, world *test.World) *grpc.ClientConn {
	t.Helper()

	conn, err := world.NewGRPC()
	require.NoError(t, err)

	return conn
}

func requireObservedHealth(t *testing.T, server *server.Server, service string, healthy bool) {
	t.Helper()

	observer, err := server.Observer(service, "grpc")
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		if healthy {
			return observer.Error() == nil
		}

		return observer.Error() != nil
	}, time.Second.Duration(), (10 * time.Millisecond).Duration())
}

func newGRPCHealthWorld(t *testing.T, url string, opts ...test.WorldOption) *test.World {
	t.Helper()

	opts = append(opts, test.WithWorldGRPCHealth(test.Name.String(), url, test.HealthObserve("grpc", "http")))

	return test.NewStartedWorld(t, opts...)
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
