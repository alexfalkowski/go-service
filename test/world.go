package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/alexfalkowski/go-service/cache/redis"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/database/sql/pg"
	sm "github.com/alexfalkowski/go-service/database/sql/telemetry/metrics"
	"github.com/alexfalkowski/go-service/hooks"
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

//nolint:gochecknoinits
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
	generator   token.Generator
	verfier     token.Verifier
	rt          http.RoundTripper
	limiter     *limiter.Config
	redis       *redis.Config
	telemetry   string
	secure      bool
	rest        bool
	compression bool
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
func WithWorldRedisConfig(config *redis.Config) WorldOption {
	return worldOptionFunc(func(o *worldOpts) {
		o.redis = config
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
	t *testing.T
	*fxtest.Lifecycle
	*http.ServeMux
	*Server
	*Client
	*mvc.Router
	*events.Event
	*eh.Receiver
	*Cache
	Sender client.Client
	Rest   *resty.Client
}

// NewWorld for test.
func NewWorld(t *testing.T, opts ...WorldOption) *World {
	t.Helper()

	mux := http.NewServeMux()
	lc := fxtest.NewLifecycle(t)
	logger := NewLogger(lc)
	tracer := NewOTLPTracerConfig()
	os := options(opts...)
	tranConfig := transportConfig(os)
	tlsConfig := tlsConfig(os)
	meter := meter(lc, mux, os)
	limiter := serverLimiter(lc, os)

	server := &Server{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: tranConfig,
		Meter: meter, Mux: mux, Limiter: limiter,
		Verifier: os.verfier, VerifyAuth: os.verfier != nil,
	}
	server.Register()

	client := &Client{
		Lifecycle: lc, Logger: logger, Tracer: tracer, Transport: tranConfig,
		Meter: meter, TLS: tlsConfig, Generator: os.generator,
		Compression: os.compression, RoundTripper: os.rt,
	}

	views := mvc.NewViews(mvc.ViewsParams{FS: &Views, Patterns: mvc.Patterns{"views/*.tmpl"}})
	router := mvc.NewRouter(mux, views)

	rest.Register(mux, Content)
	rpc.Register(mux, Content, Pool)
	pg.Register(client.NewTracer(), logger)

	restClient := restClient(client, os)

	h, err := hooks.New(NewHook())
	runtime.Must(err)

	receiver := eh.NewReceiver(mux, hh.NewWebhook(h))

	sender, err := eh.NewSender(hh.NewWebhook(h), eh.WithSenderRoundTripper(os.rt))
	runtime.Must(err)

	cache := redisCache(lc, client.Logger, server.Meter, os)

	return &World{
		t:         t,
		Lifecycle: lc, ServeMux: mux,
		Server: server, Client: client,
		Router: router, Rest: restClient,
		Receiver: receiver, Sender: sender,
		Cache: cache,
	}
}

// RegisterEvents for world.
func (w *World) RegisterEvents(ctx context.Context) {
	w.Receiver.Register(ctx, "/events", func(_ context.Context, e events.Event) { w.Event = &e })
}

// EventsContext for world.
func (w *World) EventsContext(ctx context.Context) context.Context {
	addr := w.Server.Transport.HTTP.Address

	return events.ContextWithTarget(ctx, fmt.Sprintf("http://%s/events", addr))
}

// Start the world.
func (w *World) Start() {
	w.Lifecycle.RequireStart()
}

// Stop the world.
func (w *World) Stop() {
	w.Lifecycle.RequireStop()
}

// Request for the world.
func (w *World) Request(ctx context.Context, protocol, method, path string, header http.Header, body io.Reader) (*http.Response, string, error) {
	client := w.Client.NewHTTP()
	addr := w.Server.Transport.HTTP.Address

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s://%s/%s", protocol, addr, path), body)
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

// OpenDatabase for world.
func (w *World) OpenDatabase() *mssqlx.DBs {
	dbs, err := pg.Open(pg.OpenParams{Lifecycle: w.Lifecycle, Config: NewPGConfig()})
	runtime.Must(err)

	sm.Register(dbs, w.Server.Meter)

	return dbs
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

func redisCache(lc fx.Lifecycle, logger *zap.Logger, meter metric.Meter, os *worldOpts) *Cache {
	if os.redis == nil {
		os.redis = NewRedisConfig("redis", "snappy", "proto")
	}

	return &Cache{
		Lifecycle: lc,
		Redis:     os.redis,
		Logger:    logger,
		Meter:     meter,
	}
}
