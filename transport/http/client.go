package http

import (
	"crypto/tls"
	"net/http"
	"time"

	st "github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	nh "github.com/alexfalkowski/go-service/net/http"
	v1 "github.com/alexfalkowski/go-service/net/http/v1"
	v2 "github.com/alexfalkowski/go-service/net/http/v2"
	r "github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security/token"
	t "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	tkn "github.com/alexfalkowski/go-service/transport/http/security/token"
	logger "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	hm "github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	ht "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"github.com/klauspost/compress/gzhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(opts *clientOpts) }

var none = clientOptionFunc(func(_ *clientOpts) {})

type clientOpts struct {
	tracer       trace.Tracer
	meter        metric.Meter
	roundTripper http.RoundTripper
	gen          token.Generator
	logger       *zap.Logger
	retry        *r.Config
	tls          *tls.Config
	userAgent    env.UserAgent
	timeout      time.Duration
	breaker      bool
	compression  bool
	version      nh.Version
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) { f(o) }

// WithClientBreaker for HTTP.
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
		o.timeout = t.MustParseDuration(timeout)
	})
}

// WithClientRoundTripper for HTTP.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
func WithClientRetry(cfg *r.Config) ClientOption {
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

// WithUserAgent for HTTP.
func WithClientUserAgent(userAgent env.UserAgent) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.userAgent = userAgent
	})
}

// WithClientTLS for HTTP.
func WithClientTLS(sec *st.Config) (ClientOption, error) {
	if !st.IsEnabled(sec) {
		return none, nil
	}

	conf, err := st.NewConfig(sec)
	if err != nil {
		return none, err
	}

	opt := clientOptionFunc(func(o *clientOpts) {
		o.tls = conf
	})

	return opt, nil
}

// WithClientVersion for HTTP.
func WithClientVersion(version nh.Version) ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.version = version
	})
}

// NewRoundTripper for HTTP.
func NewRoundTripper(opts ...ClientOption) http.RoundTripper {
	os := &clientOpts{tracer: noop.Tracer{}}
	for _, o := range opts {
		o.apply(os)
	}

	hrt := roundTripper(os)

	if os.compression {
		hrt = gzhttp.Transport(hrt, gzhttp.TransportEnableGzip(true))
	}

	if os.retry != nil {
		hrt = retry.NewRoundTripper(os.retry, hrt)
	}

	if os.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	if os.gen != nil {
		hrt = tkn.NewRoundTripper(os.gen, hrt)
	}

	if os.logger != nil {
		hrt = logger.NewRoundTripper(os.logger, hrt)
	}

	if os.meter != nil {
		rt := hm.NewRoundTripper(os.meter, hrt)

		hrt = rt
	}

	hrt = ht.NewRoundTripper(os.tracer, hrt)
	hrt = meta.NewRoundTripper(os.userAgent, hrt)

	return hrt
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) *http.Client {
	os := &clientOpts{tracer: noop.Tracer{}}
	for _, o := range opts {
		o.apply(os)
	}

	client := &http.Client{
		Transport: NewRoundTripper(opts...),
		Timeout:   os.timeout,
	}

	return client
}

func roundTripper(os *clientOpts) http.RoundTripper {
	hrt := os.roundTripper
	if hrt != nil {
		return hrt
	}

	if os.version == 0 {
		os.version = nh.V1
	}

	switch os.version {
	case nh.V1:
		hrt = v1.Transport(os.tls)
	case nh.V2:
		hrt = v2.Transport(os.tls)
	}

	return hrt
}
