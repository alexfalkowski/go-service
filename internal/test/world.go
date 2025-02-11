package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alexfalkowski/go-service/cache"
	"github.com/alexfalkowski/go-service/cache/cachego"
	cc "github.com/alexfalkowski/go-service/cache/config"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	sm "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/debug"
	"github.com/alexfalkowski/go-service/hooks"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/limiter"
	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	eh "github.com/alexfalkowski/go-service/transport/events/http"
	ht "github.com/alexfalkowski/go-service/transport/http"
	hh "github.com/alexfalkowski/go-service/transport/http/hooks"
	"github.com/alexfalkowski/go-service/transport/meta"
	events "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/go-resty/resty/v2"
	"github.com/linxGnu/mssqlx"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
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
	verfier     token.Verifier
	rt          http.RoundTripper
	generator   token.Generator
	logger      *zap.Logger
	limiter     *limiter.Config
	pg          *pg.Config
	telemetry   string
	secure      bool
	compression bool
	http        bool
	grpc        bool
	debug       bool
	rest        bool
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
		o.verfier = verifier
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
func WithWorldLogger(logger *zap.Logger) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.logger = logger
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
	*zap.Logger
	Tracer *tracer.Config
	PG     *pg.Config
	*Server
	*Client
	*events.Event
	*eh.Receiver
	cc.Cache
	Sender client.Client
	Rest   *resty.Client
}

// NewWorld for test.
func NewWorld(t fxtest.TB, opts ...WorldOption) *World {
	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(t)
	tracer := NewOTLPTracerConfig()
	id := id.Default
	os := options(opts...)

	if os.logger == nil {
		os.logger = NewLogger(lc)
	}

	tranConfig := transportConfig(os)
	debugConfig := debugConfig(os)
	tlsConfig := tlsConfig(os)
	meter := meter(lc, mux, os)
	limiter := serverLimiter(lc, os)
	pgConfig := pgConfig(os)

	server := &Server{
		Lifecycle: lc, Logger: os.logger, Tracer: tracer,
		TransportConfig: tranConfig, DebugConfig: debugConfig,
		Meter: meter, Mux: mux, Limiter: limiter,
		Verifier: os.verfier, VerifyAuth: os.verfier != nil, ID: id,
		RegisterHTTP: os.http, RegisterGRPC: os.grpc, RegisterDebug: os.debug,
	}
	server.Register()

	client := &Client{
		Lifecycle: lc, Logger: os.logger, Tracer: tracer, Transport: tranConfig,
		Meter: meter, TLS: tlsConfig, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
	}

	views := mvc.NewViews(mvc.ViewsParams{FS: &Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
	mvc.Register(mux, views)

	restClient := restClient(client, os)

	h, err := hooks.New(NewHook())
	runtime.Must(err)

	receiver := eh.NewReceiver(mux, hh.NewWebhook(h, id))

	sender, err := eh.NewSender(hh.NewWebhook(h, id), eh.WithSenderRoundTripper(os.rt))
	runtime.Must(err)

	cache := redisCache(lc, os.logger, meter, tracer)

	return &World{
		Logger: os.logger, Tracer: tracer,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Rest:     restClient,
		Receiver: receiver, Sender: sender,
		Cache: cache, PG: pgConfig,
	}
}

func (w *World) Register() {
	rest.Register(w.ServeMux, Content)
	rpc.Register(w.ServeMux, Content, Pool)
	pg.Register(w.NewTracer(), w.Logger)
}

// ServerHost for world.
func (w *World) ServerHost() string {
	return w.Server.TransportConfig.HTTP.Address
}

// DebugHost for world.
func (w *World) DebugHost() string {
	return w.Server.DebugConfig.Address
}

// RegisterEvents for world.
func (w *World) RegisterEvents(ctx context.Context) {
	w.Receiver.Register(ctx, "/events", func(_ context.Context, e events.Event) { w.Event = &e })
}

// EventsContext for world.
func (w *World) EventsContext(ctx context.Context) context.Context {
	return events.ContextWithTarget(ctx, fmt.Sprintf("http://%s/events", w.ServerHost()))
}

// ResponseWithBody for the world.
func (w *World) ResponseWithBody(ctx context.Context, protocol, address, method, path string, header http.Header, body io.Reader) (*http.Response, string, error) {
	client := w.Client.NewHTTP()

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

	return res, strings.TrimSpace(string(data)), nil
}

// HTTPResponseNoBody for the world.
func (w *World) ResponseWithNoBody(ctx context.Context, protocol, address, method, path string, header http.Header, body io.Reader) (*http.Response, error) {
	client := w.Client.NewHTTP()

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
	dbs, err := pg.Open(pg.OpenParams{Lifecycle: w.Lifecycle, Config: w.PG})
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

func meter(lc fx.Lifecycle, mux *http.ServeMux, os *worldOpts) metric.Meter {
	if os.telemetry == "otlp" {
		return NewOTLPMeter(lc)
	}

	config := NewPrometheusMetricsConfig()
	ht.RegisterMetrics(config, mux)

	return NewMeter(lc, config)
}

func restClient(client *Client, os *worldOpts) *resty.Client {
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

func redisCache(lc fx.Lifecycle, logger *zap.Logger, meter metric.Meter, tracer *tracer.Config) cc.Cache {
	cfg := NewCacheConfig("redis", "snappy", "json", "redis")

	cachego, err := cachego.New(cfg)
	runtime.Must(err)

	params := cache.Params{
		Lifecycle:  lc,
		Config:     cfg,
		Compressor: Compressor,
		Encoder:    Encoder,
		Pool:       Pool,
		Cache:      cachego,
		Tracer:     NewTracer(lc, tracer, logger),
		Logger:     logger,
		Meter:      meter,
	}

	cache, err := cache.New(params)
	runtime.Must(err)

	return cache
}

func pgConfig(os *worldOpts) *pg.Config {
	if os.pg != nil {
		return os.pg
	}

	return NewPGConfig()
}
