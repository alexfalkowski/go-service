package health_test

import (
	"net/http/httptest"
	"testing"

	"github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/grpc/meta"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/time"
	grpchealth "github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"github.com/alexfalkowski/go-sync"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
	)
	requireGRPCReady(t, world)

	ctx := meta.WithAttributes(t.Context(),
		meta.WithRequestID(meta.String("test-id")),
		meta.WithUserAgent(meta.String("test-user-agent")),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, health.Serving, resp.GetStatus())
}

func TestInvalidCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)
	requireObservedHealth(t, world.GRPCHealth, test.Name.String(), false)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	md := meta.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := meta.NewOutgoingContext(t.Context(), md)

	resp, err := client.Check(ctx, req)
	require.NoError(t, err)

	require.Equal(t, health.NotServing, resp.GetStatus())
}

func TestNotFoundCheck(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.Request{Service: "bob"}

	md := meta.New(map[string]string{"request-id": "test-id", "user-agent": "test-user-agent"})
	ctx := meta.NewOutgoingContext(t.Context(), md)

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

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	resp, err := client.Check(t.Context(), req)
	require.NoError(t, err)

	require.Equal(t, health.Serving, resp.GetStatus())
}

func TestList(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("200"),
		test.WithWorldTelemetry("otlp"),
		test.WithWorldServerLimiter(test.NewLimiterConfig("user-agent", "1s", 100)),
	)
	requireGRPCReady(t, world)

	ctx := meta.WithAttributes(t.Context(),
		meta.WithRequestID(meta.String("test-id")),
		meta.WithUserAgent(meta.String("test-user-agent")),
	)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.ListRequest{}

	resp, err := client.List(ctx, req)
	require.NoError(t, err)

	expected := map[string]*health.Response{
		test.Name.String(): {Status: health.Serving},
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

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, health.Serving, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)
}

func TestInvalidWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)
	requireObservedHealth(t, world.GRPCHealth, test.Name.String(), false)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, health.NotServing, resp.GetStatus())
	requireWatchStaysOpen(t, cancel, wc)
}

func TestNotFoundWatch(t *testing.T) {
	world := newGRPCHealthWorld(t, test.StatusURL("500"), test.WithWorldTelemetry("otlp"))
	requireGRPCReady(t, world)

	conn := requireGRPCConn(t, world)
	defer conn.Close()

	client := health.NewClient(conn)
	req := &health.Request{Service: "bob"}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)
	require.Equal(t, health.ServiceUnknown, resp.GetStatus())
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

	client := health.NewClient(conn)
	req := &health.Request{Service: test.Name.String()}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	wc, err := client.Watch(ctx, req)
	require.NoError(t, err)

	resp, err := wc.Recv()
	require.NoError(t, err)

	require.Equal(t, health.Serving, resp.GetStatus())
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

	watcher := grpchealth.NewServer(grpchealth.ServerParams{Server: world.GRPCHealth})
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	stream := newWatchStream(ctx)
	errCh := make(chan error, 1)
	go func() {
		errCh <- watcher.Watch(&health.Request{Service: test.Name.String()}, stream)
	}()

	resp := requireWatchResponse(t, stream.responses)
	require.Equal(t, health.Serving, resp.GetStatus())

	unhealthy.Store(true)

	require.Eventually(t, func() bool {
		select {
		case resp = <-stream.responses:
			return resp.GetStatus() == health.NotServing
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

func requireWatchStaysOpen(t *testing.T, cancel context.CancelFunc, wc health.WatchClient) {
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

func requireWatchResponse(t *testing.T, responses <-chan *health.Response) *health.Response {
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
	return &watchStream{ctx: ctx, responses: make(chan *health.Response, 4)}
}

type watchStream struct {
	grpc.ServerStream
	ctx       context.Context
	responses chan *health.Response
}

func (w *watchStream) Context() context.Context {
	return w.ctx
}

func (w *watchStream) Send(resp *health.Response) error {
	select {
	case <-w.ctx.Done():
		return w.ctx.Err()
	case w.responses <- resp:
		return nil
	}
}

func (*watchStream) SetHeader(meta.Map) error {
	return nil
}

func (*watchStream) SendHeader(meta.Map) error {
	return nil
}

func (*watchStream) SetTrailer(meta.Map) {}

func (*watchStream) SendMsg(any) error {
	return nil
}

func (*watchStream) RecvMsg(any) error {
	return nil
}
