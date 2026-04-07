package test

import (
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
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
	sdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/linxGnu/mssqlx"
	"go.uber.org/fx/fxtest"
)

func init() {
	telemetry.Register()
	grpc.Register(FS)
	th.Register(FS)
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
func NewWorld(tb testing.TB, opts ...WorldOption) *World {
	tb.Helper()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(tb)
	tracer := NewOTLPTracerConfig()
	generator := uuid.NewGenerator()
	os := worldOptions(opts...)

	logger := createLogger(lc, os)
	transportCfg := transportConfig(os)
	debugCfg := debugConfig(os)
	tlsCfg := tlsConfig(os)
	meter := meter(lc, mux, os)
	server := &Server{
		Lifecycle: lc, Logger: logger, Tracer: tracer,
		TransportConfig: transportCfg, DebugConfig: debugCfg,
		Meter: meter, Mux: mux,
		GRPCLimiter: NewGRPCServerLimiter(lc, LimiterKeyMap, os.serverLimiter),
		HTTPLimiter: NewHTTPServerLimiter(lc, LimiterKeyMap, os.serverLimiter),
		Verifier:    os.verifier, Generator: generator,
		RegisterHTTP: os.http, RegisterGRPC: os.grpc, RegisterDebug: os.debug,
	}
	server.Register()

	client := &Client{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: transportCfg,
		Meter: meter, TLS: tlsCfg, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
		HTTPLimiter: NewHTTPClientLimiter(lc, LimiterKeyMap, os.clientLimiter),
		GRPCLimiter: NewGRPCClientLimiter(lc, LimiterKeyMap, os.clientLimiter),
	}

	registerMVC(mux, logger.Logger)
	registerRest(mux)

	receiver, sender := NewEvents(mux, os.rt, generator)

	world := &World{
		t:      tb,
		Logger: logger, Tracer: tracer,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Rest:     restClient(client, os),
		Receiver: receiver, Sender: sender,
		Cache: newWorldCache(tb, lc, os), PG: os.pg,
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
	*sdk.Event
	*events.Receiver
	*cache.Cache
	Sender client.Client
	Rest   *rest.Client

	DB         *mssqlx.DBs
	HTTPHealth *health.Server
	GRPCHealth *health.Server
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

// ResponseWithBody issues an HTTP request through the world's client and returns the response plus a trimmed body string.
func (w *World) ResponseWithBody(ctx context.Context, url, method string, header http.Header, body io.Reader) (*http.Response, string, error) {
	client := w.NewHTTP()

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	runtime.Must(err)

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, strings.Empty, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		runtime.Must(err)
	}

	return res, bytes.String(bytes.TrimSpace(data)), res.Body.Close()
}

// ResponseWithNoBody issues an HTTP request through the world's client and closes the response body before returning.
func (w *World) ResponseWithNoBody(ctx context.Context, url, method string, header http.Header) (*http.Response, error) {
	client := w.NewHTTP()

	req, err := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	runtime.Must(err)

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, res.Body.Close()
}

func (w *World) namedURL(host, path string) string {
	url, err := url.JoinPath(host, Name.String(), path)
	runtime.Must(err)

	return url
}

func (w *World) pathURL(host, path string) string {
	url, err := url.JoinPath(host, path)
	runtime.Must(err)

	return url
}

func (w *World) url(protocol, address string) string {
	_, host, _ := net.SplitNetworkAddress(address)

	return strings.Concat(protocol, "://", host)
}
