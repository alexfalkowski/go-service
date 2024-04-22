package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net"
	r "github.com/alexfalkowski/go-service/retry"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/http/breaker"
	"github.com/alexfalkowski/go-service/transport/http/meta"
	"github.com/alexfalkowski/go-service/transport/http/retry"
	lzap "github.com/alexfalkowski/go-service/transport/http/telemetry/logger/zap"
	hm "github.com/alexfalkowski/go-service/transport/http/telemetry/metrics"
	ht "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ClientOption for HTTP.
type ClientOption interface{ apply(opts *clientOptions) }

type clientOptions struct {
	logger       *zap.Logger
	tracer       trace.Tracer
	meter        metric.Meter
	retry        *r.Config
	userAgent    string
	breaker      bool
	roundTripper http.RoundTripper
}

type clientOptionFunc func(*clientOptions)

func (f clientOptionFunc) apply(o *clientOptions) { f(o) }

// WithClientRoundTripper for HTTP.
func WithClientRoundTripper(rt http.RoundTripper) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.roundTripper = rt
	})
}

// WithClientRetry for HTTP.
func WithClientRetry(cfg *r.Config) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.retry = cfg
	})
}

// WithClientBreaker for HTTP.
func WithClientBreaker() ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.breaker = true
	})
}

// WithClientLogger for HTTP.
func WithClientLogger(logger *zap.Logger) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.logger = logger
	})
}

// WithClientTracer for HTTP.
func WithClientTracer(tracer trace.Tracer) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.tracer = tracer
	})
}

// WithClientMetrics for HTTP.
func WithClientMetrics(meter metric.Meter) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.meter = meter
	})
}

// WithUserAgent for HTTP.
func WithClientUserAgent(userAgent string) ClientOption {
	return clientOptionFunc(func(o *clientOptions) {
		o.userAgent = userAgent
	})
}

// NewRoundTripper for HTTP.
func NewRoundTripper(opts ...ClientOption) (http.RoundTripper, error) {
	os := &clientOptions{tracer: tracer.NewNoopTracer()}
	for _, o := range opts {
		o.apply(os)
	}

	hrt := os.roundTripper
	if hrt == nil {
		hrt = Transport()
	}

	if os.retry != nil {
		hrt = retry.NewRoundTripper(os.retry, hrt)
	}

	if os.breaker {
		hrt = breaker.NewRoundTripper(hrt)
	}

	if os.logger != nil {
		hrt = lzap.NewRoundTripper(os.logger, hrt)
	}

	if os.meter != nil {
		rt, err := hm.NewRoundTripper(os.meter, hrt)
		if err != nil {
			return nil, err
		}

		hrt = rt
	}

	hrt = ht.NewRoundTripper(os.tracer, hrt)
	hrt = meta.NewRoundTripper(os.userAgent, hrt)

	return hrt, nil
}

// NewClient for HTTP.
func NewClient(opts ...ClientOption) (*http.Client, error) {
	defaultOptions := &clientOptions{tracer: tracer.NewNoopTracer()}
	for _, o := range opts {
		o.apply(defaultOptions)
	}

	rt, err := NewRoundTripper(opts...)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: rt,
		Timeout:   time.Timeout,
	}

	return client, nil
}

// Transport for HTTP.
func Transport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.DialContext = net.DialContext

	return t
}
