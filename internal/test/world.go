package test

import (
	"io"
	"net/url"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/id/uuid"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/telemetry"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	tg "github.com/alexfalkowski/go-service/v2/transport/grpc"
	th "github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
	sdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"go.uber.org/fx/fxtest"
)

func init() {
	telemetry.Register()
	tg.Register(FS)
	th.Register(FS)
	Encoder.Register("error", NewEncoder(ErrFailed))
	Compressor.Register("error", NewCompressor(ErrFailed))
}

// WorldOption for test.
type WorldOption interface {
	apply(opts *worldOpts)
}

type worldOpts struct {
	verifier      token.Verifier
	rt            http.RoundTripper
	generator     token.Generator
	logger        *logger.Logger
	clientLimiter *limiter.Config
	serverLimiter *limiter.Config
	pg            *pg.Config
	telemetry     string
	loggerConfig  string
	secure        bool
	compression   bool
	http          bool
	grpc          bool
	debug         bool
	rest          bool
}

type worldOptionFunc func(*worldOpts)

func (f worldOptionFunc) apply(o *worldOpts) {
	f(o)
}

// WithWorldSecure for test.
func WithWorldSecure() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.secure = true
	})
}

// WithWorldTelemetry for test.
func WithWorldTelemetry(kind string) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.telemetry = kind
	})
}

// WithWorldClientLimiter for test.
func WithWorldClientLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.clientLimiter = config
	})
}

// WithWorldServerLimiter for test.
func WithWorldServerLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.serverLimiter = config
	})
}

// WithWorldCompression for test.
func WithWorldCompression() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.compression = true
	})
}

// WithWorldRoundTripper for test.
func WithWorldRoundTripper(rt http.RoundTripper) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.rt = rt
	})
}

// WithWorldHTTP for test.
func WithWorldHTTP() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.http = true
	})
}

// WithWorldGRPC for test.
func WithWorldGRPC() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.grpc = true
	})
}

// WithWorldDebug for test.
func WithWorldDebug() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.debug = true
	})
}

func worldOptions(opts ...WorldOption) *worldOpts {
	os := &worldOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	return os
}

// NewWorld for test.
func NewWorld(t fxtest.TB, opts ...WorldOption) *World {
	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(t)
	tracer := NewOTLPTracerConfig()
	generator := uuid.NewGenerator()
	os := worldOptions(opts...)

	logger := createLogger(lc, os)
	tranConfig := transportConfig(os)
	debugConfig := debugConfig(os)
	tlsConfig := tlsConfig(os)
	meter := meter(lc, mux, os)
	pgConfig := pgConfig(os)

	server := &Server{
		Lifecycle: lc, Logger: logger, Tracer: tracer,
		TransportConfig: tranConfig, DebugConfig: debugConfig,
		Meter: meter, Mux: mux,
		GRPCLimiter: NewGRPCServerLimiter(lc, LimiterKeyMap, os.serverLimiter),
		HTTPLimiter: NewHTTPServerLimiter(lc, LimiterKeyMap, os.serverLimiter),
		Verifier:    os.verifier, Generator: generator,
		RegisterHTTP: os.http, RegisterGRPC: os.grpc, RegisterDebug: os.debug,
	}
	server.Register()

	client := &Client{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: tranConfig,
		Meter: meter, TLS: tlsConfig, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
		HTTPLimiter: NewHTTPClientLimiter(lc, LimiterKeyMap, os.clientLimiter),
		GRPCLimiter: NewGRPCClientLimiter(lc, LimiterKeyMap, os.clientLimiter),
	}

	registerMVC(mux, logger.Logger)
	registerRest(mux)

	receiver, sender := NewEvents(mux, os.rt, generator)

	return &World{
		Logger: logger, Tracer: tracer,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Rest:     restClient(client, os),
		Receiver: receiver, Sender: sender,
		Cache: redisCache(lc), PG: pgConfig,
	}
}

// World for test.
type World struct {
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
}

// Register all packages.
func (w *World) Register() {
	w.registerRPC()
	w.registerDatabase()
	w.registerTelemetry()
}

// HandleHello for world.
func (w *World) HandleHello() {
	w.HandleFunc("GET /hello", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(strings.Bytes("hello!"))
	})
}

// NamedServerURL for world.
func (w *World) NamedServerURL(protocol, path string) string {
	return w.namedURL(w.ServerURL(protocol), path)
}

// PathServerURL for world.
func (w *World) PathServerURL(protocol, path string) string {
	return w.pathURL(w.ServerURL(protocol), path)
}

// ServerURL for world.
func (w *World) ServerURL(protocol string) string {
	return w.url(protocol, w.TransportConfig.HTTP.Address)
}

// NamedDebugURL for world.
func (w *World) NamedDebugURL(protocol, path string) string {
	return w.namedURL(w.DebugURL(protocol), path)
}

// PathDebugURL for world.
func (w *World) PathDebugURL(protocol, path string) string {
	return w.pathURL(w.DebugURL(protocol), path)
}

// DebugURL for world.
func (w *World) DebugURL(protocol string) string {
	return w.url(protocol, w.DebugConfig.Address)
}

// ResponseWithBody for the world.
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

// ResponseWithNoBody for the world.
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

func transportConfig(os *worldOpts) *transport.Config {
	if os.secure {
		return NewSecureTransportConfig()
	}

	return NewInsecureTransportConfig()
}

func debugConfig(os *worldOpts) *debug.Config {
	if os.secure {
		return NewSecureDebugConfig()
	}

	return NewInsecureDebugConfig()
}

func tlsConfig(os *worldOpts) *tls.Config {
	if os.secure {
		return NewTLSClientConfig()
	}

	return nil
}
