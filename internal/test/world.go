package test

import (
	"context"
	"fmt"
	"io"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	eh "github.com/alexfalkowski/go-service/v2/transport/http/events"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func init() {
	meta.RegisterKeys()
	tracer.Register()
	Encoder.Register("error", NewEncoder(ErrFailed))
	Compressor.Register("error", NewCompressor(ErrFailed))
}

// WorldOption for test.
type WorldOption interface {
	apply(opts *worldOpts)
}

type worldOpts struct {
	verifier     token.Verifier
	rt           http.RoundTripper
	generator    token.Generator
	logger       *logger.Logger
	limiter      *limiter.Config
	pg           *pg.Config
	telemetry    string
	loggerConfig string
	secure       bool
	compression  bool
	http         bool
	grpc         bool
	debug        bool
	rest         bool
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

// WithWorldLimiter for test.
func WithWorldLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.limiter = config
	})
}

// WithWorldLimiter for test.
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

func options(opts ...WorldOption) *worldOpts {
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
	id := &id.UUID{}
	os := options(opts...)

	logger := createLogger(lc, os)
	tranConfig := transportConfig(os)
	debugConfig := debugConfig(os)
	tlsConfig := tlsConfig(os)
	meter := meter(lc, mux, os)
	limiter := serverLimiter(lc, os)
	pgConfig := pgConfig(os)

	server := &Server{
		Lifecycle: lc, Logger: logger, Tracer: tracer,
		TransportConfig: tranConfig, DebugConfig: debugConfig,
		Meter: meter, Mux: mux, Limiter: limiter,
		Verifier: os.verifier, ID: id,
		RegisterHTTP: os.http, RegisterGRPC: os.grpc, RegisterDebug: os.debug,
	}
	server.Register()

	client := &Client{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: tranConfig,
		Meter: meter, TLS: tlsConfig, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
	}

	registerMVC(mux, logger.Logger)
	registerRest(mux)

	receiver, sender := NewEvents(mux, os.rt, id)
	cache := redisCache(lc, logger, meter, tracer)

	return &World{
		Logger: logger, Tracer: tracer,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Rest:     restClient(client, os),
		Receiver: receiver, Sender: sender,
		Cache: cache, PG: pgConfig,
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
	*events.Event
	*eh.Receiver
	Cache  cacher.Cache
	Sender client.Client
	Rest   *rest.Client
}

// Register all packages.
func (w *World) Register() {
	w.registerRPC()
	w.registerDatabase()
	w.registerTelemetry()
}

// InsecureServerHost for world.
func (w *World) InsecureServerHost() string {
	return w.TransportConfig.HTTP.Address
}

// SecureServerHost for world.
func (w *World) SecureServerHost() string {
	return w.TransportConfig.HTTP.Address
}

// InsecureDebugHost for world.
func (w *World) InsecureDebugHost() string {
	return w.DebugConfig.Address
}

// SecureDebugHost for world.
func (w *World) SecureDebugHost() string {
	return w.DebugConfig.Address
}

// ResponseWithBody for the world.
func (w *World) ResponseWithBody(ctx context.Context, protocol, address, method, path string, header http.Header, body io.Reader) (*http.Response, string, error) {
	client := w.NewHTTP()

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s://%s/%s", protocol, address, path), body)
	runtime.Must(err)

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	data, err := io.ReadAll(res.Body)
	runtime.Must(err)

	return res, bytes.String(bytes.TrimSpace(data)), res.Body.Close()
}

// ResponseWithNoBody for the world.
func (w *World) ResponseWithNoBody(ctx context.Context, protocol, address, method, path string, header http.Header) (*http.Response, error) {
	client := w.NewHTTP()

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s://%s/%s", protocol, address, path), http.NoBody)
	runtime.Must(err)

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, res.Body.Close()
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

func serverLimiter(lc fx.Lifecycle, os *worldOpts) *limiter.Limiter {
	if os.limiter != nil {
		l, err := limiter.New(lc, os.limiter)
		runtime.Must(err)

		return l
	}

	return nil
}
