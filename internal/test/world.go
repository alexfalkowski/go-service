package test

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/cache"
	"github.com/alexfalkowski/go-service/v2/cache/cacher"
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/database/sql/pg"
	sm "github.com/alexfalkowski/go-service/v2/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/debug"
	"github.com/alexfalkowski/go-service/v2/hooks"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/errors"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	eh "github.com/alexfalkowski/go-service/v2/transport/http/events"
	hh "github.com/alexfalkowski/go-service/v2/transport/http/hooks"
	hm "github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/meta"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/linxGnu/mssqlx"
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

// WithWorldRest for test.
func WithWorldRest() WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.rest = true
	})
}

// WithWorldLimiter for test.
func WithWorldLimiter(config *limiter.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.limiter = config
	})
}

// WithWorldToken for test.
func WithWorldToken(generator token.Generator, verifier token.Verifier) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.generator = generator
		o.verifier = verifier
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

// WithWorldRedisConfig for test.
func WithWorldPGConfig(config *pg.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.pg = config
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

// WithWorldLogger for test.
func WithWorldLogger(logger *logger.Logger) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.logger = logger
	})
}

// WithWorldLogger for test.
func WithWorldLoggerConfig(config string) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.loggerConfig = config
	})
}

func options(opts ...WorldOption) *worldOpts {
	os := &worldOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	return os
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

	mvc.Register(mvc.RegisterParams{
		Mux:         mux,
		FunctionMap: mvc.NewFunctionMap(mvc.FunctionMapParams{Logger: logger.Logger}),
		FileSystem:  FileSystem,
		Layout:      Layout,
	})

	rest.Register(mux, Content, Pool)

	h, err := hooks.New(FS, NewHook())
	runtime.Must(err)

	receiver := eh.NewReceiver(mux, hh.NewWebhook(h, id))

	sender, err := eh.NewSender(hh.NewWebhook(h, id), eh.WithSenderRoundTripper(os.rt))
	runtime.Must(err)

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

// Register all packages.
func (w *World) Register() {
	rpc.Register(w.ServeMux, Content, Pool)
	pg.Register(w.NewTracer(), w.Logger)
	errors.Register(errors.NewHandler(w.Logger))
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

// RegisterEvents for world.
func (w *World) RegisterEvents(ctx context.Context) {
	w.Receiver.Register(ctx, "/events", func(_ context.Context, e events.Event) { w.Event = &e })
}

// EventsContext for world.
func (w *World) EventsContext(ctx context.Context) context.Context {
	return events.ContextWithTarget(ctx, fmt.Sprintf("http://%s/events", w.InsecureServerHost()))
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

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	runtime.Must(err)

	return res, bytes.String(bytes.TrimSpace(data)), nil
}

// HTTPResponseNoBody for the world.
func (w *World) ResponseWithNoBody(ctx context.Context, protocol, address, method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	client := w.NewHTTP()

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s://%s/%s", protocol, address, path), body)
	runtime.Must(err)

	req.Header = header

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return res, nil
}

// OpenDatabase for world.
func (w *World) OpenDatabase() (*mssqlx.DBs, error) {
	dbs, err := pg.Open(w.Lifecycle, FS, w.PG)
	if err != nil {
		return nil, err
	}

	sm.Register(dbs, w.Server.Meter)

	return dbs, err
}

// RegisterHandlers for test.
func RegisterHandlers[Res any](path string, h content.Handler[Res]) {
	rest.Delete(path, h)
	rest.Get(path, h)
}

// RegisterRequestHandlers for test.
func RegisterRequestHandlers[Req any, Res any](path string, h content.RequestHandler[Req, Res]) {
	rest.Post(path, h)
	rest.Put(path, h)
	rest.Patch(path, h)
}

func createLogger(lc fx.Lifecycle, os *worldOpts) *logger.Logger {
	if os.logger != nil {
		return os.logger
	}

	var config *logger.Config

	switch os.loggerConfig {
	case "json":
		config = NewJSONLoggerConfig()
	case "text":
		config = NewTextLoggerConfig()
	case "tilt":
		config = NewTintLoggerConfig()
	case "otlp":
		config = NewOTLPLoggerConfig()
	default:
		config = NewOTLPLoggerConfig()
	}

	return NewLogger(lc, config)
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

func meter(lc fx.Lifecycle, mux *http.ServeMux, os *worldOpts) *metrics.Meter {
	if os.telemetry == "otlp" {
		return NewOTLPMeter(lc)
	}

	config := NewPrometheusMetricsConfig()
	hm.Register(config, mux)

	return NewMeter(lc, config)
}

func restClient(client *Client, os *worldOpts) *rest.Client {
	if os.rest {
		return rest.NewClient(
			rest.WithClientRoundTripper(client.NewHTTP().Transport),
			rest.WithClientTimeout("10s"),
		)
	}

	return rest.NewClient()
}

func serverLimiter(lc fx.Lifecycle, os *worldOpts) *limiter.Limiter {
	if os.limiter != nil {
		l, err := limiter.New(lc, os.limiter)
		runtime.Must(err)

		return l
	}

	return nil
}

func redisCache(lc fx.Lifecycle, logger *logger.Logger, meter *metrics.Meter, tracer *tracer.Config) cacher.Cache {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	driver, err := driver.New(FS, cfg)
	runtime.Must(err)

	params := cache.Params{
		Lifecycle:  lc,
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Driver:     driver,
		Tracer:     NewTracer(lc, tracer),
		Logger:     logger,
		Meter:      meter,
	}

	return cache.NewCache(params)
}

func pgConfig(os *worldOpts) *pg.Config {
	if os.pg != nil {
		return os.pg
	}

	return NewPGConfig()
}
