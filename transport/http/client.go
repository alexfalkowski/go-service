package http

import (
	"net/http"
	"net/url"
	"time"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/id"
	nh "github.com/alexfalkowski/go-service/net/http"
	sr "github.com/alexfalkowski/go-service/retry"
	st "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	logger "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	hm "github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	tt "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	ht "github.com/alexfalkowski/go-service/transport/http/token"
	"github.com/klauspost/compress/gzhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface {
	apply(opts *clientOpts)
}

type clientOpts struct {
	tracer       trace.Tracer
	meter        metric.Meter
	roundTripper http.RoundTripper
	gen          token.Generator
	logger       *zap.Logger
	retry        *sr.Config
	tls          *tls.Config
	id           id.Generator
	userAgent    env.UserAgent
	timeout      time.Duration
	breaker      bool
	compression  bool
	proxy        string
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

// WithClientCompression for HTTP.
func WithClientCompression() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.compression = true
	})
}

// WithClientTokenGenerator for HTTP.
func WithClientTokenGenerator(gen token.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.gen = gen
	})
}

// WithClientTimeout for HTTP.
func WithClientTimeout(timeout string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.timeout = st.MustParseDuration(timeout)
	})
}

// WithClientRoundTripper for HTTP.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
func WithClientRetry(cfg *sr.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.retry = cfg
	})
}

// WithClientBreaker for HTTP.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.breaker = true
	})
}

// WithClientLogger for HTTP.
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.logger = logger
	})
}

// WithClientTracer for HTTP.
func WithClientTracer(tracer trace.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tracer = tracer
	})
}

// WithClientMetrics for HTTP.
func WithClientMetrics(meter metric.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.meter = meter
	})
}

// WithClientUserAgent for HTTP.
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientTLS for HTTP.
func WithClientTLS(sec *tls.Config) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.tls = sec
	})
}

// WithClientID for HTTP.
func WithClientID(gen id.Generator) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.id = gen
	})
}

// WithClientID for HTTP.
func WithClientProxy(url string) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.proxy = url
	})
}

// NewRoundTripper for HTTP.
func NewRoundTripper(opts ...ClientOption) (http.RoundTripper, error) {
	os := options(opts...)

	hrt, err := roundTripper(os)
	if err != nil {
		return nil, err
	}

	if os.gen != nil {
		hrt = ht.NewRoundTripper(os.gen, hrt)
	}

	if os.compression {
		hrt = gzhttp.Transport(hrt, gzhttp.TransportEnableGzip(true))
	}

	if os.retry != nil {
		hrt = retry.NewRoundTripper(os.retry, hrt)
	}

	if os.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	if os.logger != nil {
		hrt = logger.NewRoundTripper(os.logger, hrt)
	}

	if os.meter != nil {
		hrt = hm.NewRoundTripper(os.meter, hrt)
	}

	if os.tracer != nil {
		hrt = tt.NewRoundTripper(os.tracer, hrt)
	}

	hrt = meta.NewRoundTripper(os.userAgent, os.id, hrt)

	return hrt, nil
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) (*http.Client, error) {
	os := options(opts...)

	transport, err := NewRoundTripper(opts...)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   os.timeout,
	}

	return client, nil
}

func roundTripper(os *clientOpts) (http.RoundTripper, error) {
	hrt := os.roundTripper
	if hrt != nil {
		return hrt, nil
	}

	return transport(os)
}

func transport(os *clientOpts) (*http.Transport, error) {
	var transport *http.Transport

	if tls.IsEnabled(os.tls) {
		conf, err := tls.NewConfig(os.tls)
		if err != nil {
			return nil, err
		}

		transport = nh.Transport(conf)
	} else {
		transport = nh.Transport(nil)
	}

	if os.proxy == "" {
		return transport, nil
	}

	u, err := url.Parse(os.proxy)
	if err != nil {
		return transport, err
	}

	transport.Proxy = http.ProxyURL(u)

	return transport, nil
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.timeout == 0 {
		os.timeout = 30 * time.Second
	}

	if os.id == nil {
		os.id = id.Default
	}

	return os
}
