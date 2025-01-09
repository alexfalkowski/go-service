package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	g "github.com/alexfalkowski/go-service/transport/grpc"
	h "github.com/alexfalkowski/go-service/transport/http"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Client for test.
type Client struct {
	Lifecycle    fx.Lifecycle
	Logger       *zap.Logger
	Tracer       *tracer.Config
	Transport    *transport.Config
	TLS          *tls.Config
	RoundTripper http.RoundTripper
	Meter        metric.Meter
	Generator    token.Generator
	Compression  bool
}

// NewTracer for client.
func (c *Client) NewTracer() trace.Tracer {
	tracer, err := tracer.NewTracer(c.Lifecycle, Environment, Version, Name, c.Tracer, c.Logger)
	runtime.Must(err)

	return tracer
}

// NewHTTP client for test.
func (c *Client) NewHTTP() *http.Client {
	sec, err := h.WithClientTLS(c.TLS)
	runtime.Must(err)

	tracer := c.NewTracer()
	opts := []h.ClientOption{
		h.WithClientLogger(c.Logger),
		h.WithClientRoundTripper(c.RoundTripper), h.WithClientBreaker(),
		h.WithClientTracer(tracer), h.WithClientRetry(c.Transport.HTTP.Retry),
		h.WithClientMetrics(c.Meter), h.WithClientUserAgent(UserAgent),
		h.WithClientTokenGenerator(c.Generator), h.WithClientTimeout("1m"), sec,
	}

	if c.Compression {
		opts = append(opts, h.WithClientCompression())
	}

	client := h.NewClient(opts...)

	return client
}

func (c *Client) NewGRPC() *grpc.ClientConn {
	tracer := c.NewTracer()

	sec, err := g.WithClientTLS(c.TLS)
	runtime.Must(err)

	config := &client.Config{
		Address: c.Transport.GRPC.Address,
		Retry:   c.Transport.GRPC.Retry,
	}

	opts := []g.ClientOption{
		g.WithClientUnaryInterceptors(), g.WithClientStreamInterceptors(),
		g.WithClientLogger(c.Logger), g.WithClientTracer(tracer),
		g.WithClientBreaker(), g.WithClientRetry(config.Retry),
		g.WithClientMetrics(c.Meter), g.WithClientUserAgent(UserAgent),
		g.WithClientTokenGenerator(c.Generator), g.WithClientTimeout("1m"),
		g.WithClientDialOption(), sec,
	}

	if c.Compression {
		opts = append(opts, g.WithClientCompression())
	}

	conn, err := g.NewClient(config.Address, opts...)
	runtime.Must(err)

	return conn
}
