package test

import (
	"errors"
	"io"
	"net/url"
	"testing"

	health "github.com/alexfalkowski/go-health/v2/server"
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
	v2 "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/linxGnu/mssqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"
)

func init() {
	telemetry.Register()
	grpc.Register(FS)
	http.Register(FS)
	Encoder.Register("error", NewEncoder(ErrFailed))
	Compressor.Register("error", NewCompressor(ErrFailed))
}

// NewWorld builds a shared integration test harness around the repository's transport,
// telemetry, cache, database, and event helpers.
//
// The returned World owns a fresh Fx test lifecycle, HTTP mux, server/client
// builders, telemetry config, cache handle, database config, and optional REST
// client.
//
// NewWorld completes the package-level helper registrations during construction,
// so callers only need to add any test-specific routes/handlers and then start
// the lifecycle. Most tests should prefer Start or NewStartedWorld to ensure
// cleanup is always registered with the testing framework.
//
//nolint:funlen
func NewWorld(tb testing.TB, opts ...WorldOption) *World {
	tb.Helper()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(tb)
	tracer := NewOTLPTracerConfig()
	generator := uuid.NewGenerator()
	os := worldOptions(opts...)

	logger, err := createLogger(lc, os)
	require.NoError(tb, err)
	transportCfg := transportConfig(os)
	debugCfg := debugConfig(os)
	tlsCfg := tlsConfig(os)
	meter, err := meter(lc, mux, os)
	require.NoError(tb, err)

	grpcServerLimiter, err := NewGRPCServerLimiter(lc, LimiterKeyMap, os.serverLimiter)
	require.NoError(tb, err)
	httpServerLimiter, err := NewHTTPServerLimiter(lc, LimiterKeyMap, os.serverLimiter)
	require.NoError(tb, err)

	server := &Server{
		Lifecycle: lc, Logger: logger, Tracer: tracer,
		TransportConfig: transportCfg, DebugConfig: debugCfg,
		Meter: meter, Mux: mux,
		GRPCLimiter: grpcServerLimiter,
		HTTPLimiter: httpServerLimiter,
		Verifier:    os.verifier, Generator: generator,
		RegisterHTTP: os.http, RegisterGRPC: os.grpc, RegisterDebug: os.debug,
	}
	require.NoError(tb, server.Register())

	httpClientLimiter, err := NewHTTPClientLimiter(lc, LimiterKeyMap, os.clientLimiter)
	require.NoError(tb, err)
	grpcClientLimiter, err := NewGRPCClientLimiter(lc, LimiterKeyMap, os.clientLimiter)
	require.NoError(tb, err)

	client := &Client{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: transportCfg,
		Meter: meter, TLS: tlsCfg, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
		HTTPLimiter: httpClientLimiter,
		GRPCLimiter: grpcClientLimiter,
	}
	httpClient, err := client.NewHTTP()
	require.NoError(tb, err)

	registerMVC(mux, logger.Logger)
	registerRest(mux)

	receiver, sender, err := NewEvents(mux, os.rt, generator)
	require.NoError(tb, err)

	world := &World{
		t:      tb,
		Logger: logger, Tracer: tracer,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Rest:     restClient(httpClient, os),
		Receiver: receiver, Sender: sender,
		Cache: newWorldCache(tb, lc, os), PG: os.pg,
		httpClient: httpClient,
	}

	world.registerOptions(os)
	world.registerRPC()
	world.registerDatabase()
	world.registerTelemetry()

	return world
}

// NewStartedWorld constructs a World, starts its lifecycle, and registers test cleanup.
//
// Use NewStartedWorld when the test does not need to install additional routes or
// mutate the harness before startup. For pre-start customization, call NewWorld
// and then Start.
func NewStartedWorld(tb testing.TB, opts ...WorldOption) *World {
	tb.Helper()

	world := NewWorld(tb, opts...)
	world.Start()

	return world
}

// World groups the shared components used by integration-style tests.
//
// It exposes the lifecycle, mux, logger, transport builders, cache, event
// sender/receiver, and generated configs so tests can compose realistic service
// scenarios with minimal boilerplate.
type World struct {
	t testing.TB
	*fxtest.Lifecycle
	*http.ServeMux
	*logger.Logger
	Tracer *tracer.Config
	PG     *pg.Config
	*Server
	*Client
	*v2.Event
	*events.Receiver
	*cache.Cache
	Sender client.Client
	Rest   *rest.Client

	DB         *mssqlx.DBs
	HTTPHealth *health.Server
	GRPCHealth *health.Server
	httpClient *http.Client
}

// Start starts the World's lifecycle and schedules cleanup with the test that
// created the World.
//
// Start is the preferred entry point for most tests because it pairs
// RequireStart with a cleanup-driven RequireStop.
func (w *World) Start() *World {
	w.t.Helper()
	w.t.Cleanup(func() {
		w.RequireStop()
	})

	w.RequireStart()

	return w
}

// HandleHello registers a simple HTTP hello endpoint on the world's mux.
func (w *World) HandleHello() {
	w.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(strings.Bytes("hello!"))
	})
}

// NamedServerURL returns a server URL rooted at `/<service-name>/<path>`.
func (w *World) NamedServerURL(protocol, path string) string {
	return w.namedURL(w.ServerURL(protocol), path)
}

// PathServerURL returns a server URL rooted directly at `/<path>`.
func (w *World) PathServerURL(protocol, path string) string {
	return w.pathURL(w.ServerURL(protocol), path)
}

// ServerURL returns the base URL for the world's HTTP transport using the requested scheme.
func (w *World) ServerURL(protocol string) string {
	return w.url(protocol, w.TransportConfig.HTTP.Address)
}

// NamedDebugURL returns a debug URL rooted at `/<service-name>/<path>`.
func (w *World) NamedDebugURL(protocol, path string) string {
	return w.namedURL(w.DebugURL(protocol), path)
}

// PathDebugURL returns a debug URL rooted directly at `/<path>`.
func (w *World) PathDebugURL(protocol, path string) string {
	return w.pathURL(w.DebugURL(protocol), path)
}

// DebugURL returns the base URL for the world's debug server using the requested scheme.
func (w *World) DebugURL(protocol string) string {
	return w.url(protocol, w.DebugConfig.Address)
}

// Do executes req with the world's shared HTTP client.
func (w *World) Do(req *http.Request) (*http.Response, error) {
	return w.httpClient.Do(req)
}

// GetBody issues an HTTP GET request through the world's client and returns the response plus a trimmed body string.
func (w *World) GetBody(ctx context.Context, url string, header http.Header) (*http.Response, string, error) {
	return w.ResponseWithBody(ctx, url, http.MethodGet, header, http.NoBody)
}

// GetNoBody issues an HTTP GET request through the world's client and closes the response body before returning.
func (w *World) GetNoBody(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	return w.ResponseWithNoBody(ctx, url, http.MethodGet, header)
}

// PostBody issues an HTTP POST request through the world's client and returns the response plus a trimmed body string.
func (w *World) PostBody(ctx context.Context, url string, header http.Header, body io.Reader) (*http.Response, string, error) {
	return w.ResponseWithBody(ctx, url, http.MethodPost, header, body)
}

// ResponseWithBody issues an HTTP request through the world's client and returns the response plus a trimmed body string.
func (w *World) ResponseWithBody(ctx context.Context, url, method string, header http.Header, body io.Reader) (*http.Response, string, error) {
	req, err := w.request(ctx, url, method, header, body)
	if err != nil {
		return nil, strings.Empty, err
	}

	res, err := w.Do(req)
	if err != nil {
		return nil, strings.Empty, err
	}

	data, err := w.readBody(res)
	if err != nil {
		return res, strings.Empty, err
	}

	return res, bytes.String(bytes.TrimSpace(data)), nil
}

// ResponseWithNoBody issues an HTTP request through the world's client and closes the response body before returning.
func (w *World) ResponseWithNoBody(ctx context.Context, url, method string, header http.Header) (*http.Response, error) {
	req, err := w.request(ctx, url, method, header, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := w.Do(req)
	if err != nil {
		return nil, err
	}

	return res, res.Body.Close()
}

func (w *World) namedURL(host, path string) string {
	w.t.Helper()

	url, err := url.JoinPath(host, Name.String(), path)
	require.NoError(w.t, err)
	return url
}

func (w *World) pathURL(host, path string) string {
	w.t.Helper()

	url, err := url.JoinPath(host, path)
	require.NoError(w.t, err)
	return url
}

func (w *World) url(protocol, address string) string {
	w.t.Helper()

	_, host, ok := net.SplitNetworkAddress(address)
	require.True(w.t, ok, "invalid network address: %s", address)
	return strings.Concat(protocol, "://", host)
}

func (w *World) request(ctx context.Context, url, method string, header http.Header, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header = header

	return req, nil
}

func (w *World) readBody(res *http.Response) ([]byte, error) {
	data, err := io.ReadAll(res.Body)
	return data, errors.Join(err, res.Body.Close())
}
