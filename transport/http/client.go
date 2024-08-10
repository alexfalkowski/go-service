package http

import (
	"crypto/tls"
	"net/http"
	"time"

	ct "github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/env"
	nh "github.com/alexfalkowski/go-service/net/http"
	"github.com/alexfalkowski/go-service/net/http/h2c"
	r "github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/security/token"
	t "github.com/alexfalkowski/go-service/time"
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

var none = clientOptionFunc(func(_ *clientOpts) {
})

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
	h2c          bool
}

type clientOptionFunc func(*clientOpts)

func (f clientOptionFunc) apply(o *clientOpts) {
	f(o)
}

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
func WithClientTLS(sec *ct.Config) (ClientOption, error) {
	if !ct.IsEnabled(sec) {
		return none, nil
	}

	conf, err := ct.NewConfig(sec)
	if err != nil {
		return none, err
	}

	opt := clientOptionFunc(func(o *clientOpts) {
		o.tls = conf
	})

	return opt, nil
}

// WithClientH2C for HTTP.
func WithClientH2C() ClientOption {
	return clientOptionFunc(func(o *clientOpts) {
		o.h2c = true
	})
}

// NewRoundTripper for HTTP.
func NewRoundTripper(opts ...ClientOption) http.RoundTripper {
	os := options(opts...)
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
		hrt = ht.NewRoundTripper(os.gen, hrt)
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

	hrt = meta.NewRoundTripper(os.userAgent, hrt)

	return hrt
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) *http.Client {
	os := options(opts...)
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

	if os.h2c {
		hrt = h2c.Transport()
	} else {
		hrt = nh.Transport(os.tls)
	}

	return hrt
}

func options(opts ...ClientOption) *clientOpts {
	os := &clientOpts{}
	for _, o := range opts {
		o.apply(os)
	}

	if os.timeout == 0 {
		os.timeout = 30 * time.Second
	}

	return os
}
