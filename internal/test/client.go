package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/crypto/tls"
	"github.com/alexfalkowski/go-service/id"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/telemetry/logger"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/token"
	"github.com/alexfalkowski/go-service/transport"
	g "github.com/alexfalkowski/go-service/transport/grpc"
	h "github.com/alexfalkowski/go-service/transport/http"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Client for test.
type Client struct {
	Lifecycle    fx.Lifecycle
	Logger       *logger.Logger
	Tracer       *tracer.Config
	Transport    *transport.Config
	TLS          *tls.Config
	ID           id.Generator
	RoundTripper http.RoundTripper
	Meter        metric.Meter
	Generator    token.Generator
	Compression  bool
}

// NewTracer for client.
func (c *Client) NewTracer() trace.Tracer {
	return NewTracer(c.Lifecycle, c.Tracer, c.Logger)
}

// NewHTTP client for test.
func (c *Client) NewHTTP() *http.Client {
	tracer := c.NewTracer()
	opts := []h.ClientOption{
		h.WithClientLogger(c.Logger),
		h.WithClientRoundTripper(c.RoundTripper), h.WithClientBreaker(),
		h.WithClientTracer(tracer), h.WithClientRetry(c.Transport.HTTP.Retry),
		h.WithClientMetrics(c.Meter), h.WithClientUserAgent(UserAgent),
		h.WithClientTokenGenerator(c.Generator), h.WithClientTimeout("1m"),
		h.WithClientTLS(c.TLS), h.WithClientID(c.ID),
	}

	if c.Compression {
		opts = append(opts, h.WithClientCompression())
	}

	client, err := h.NewClient(opts...)
	runtime.Must(err)

	return client
}

func (c *Client) NewGRPC() *grpc.ClientConn {
	tracer := c.NewTracer()
	opts := []g.ClientOption{
		g.WithClientUnaryInterceptors(), g.WithClientStreamInterceptors(),
		g.WithClientLogger(c.Logger), g.WithClientTracer(tracer),
		g.WithClientBreaker(), g.WithClientRetry(c.Transport.GRPC.Retry),
		g.WithClientMetrics(c.Meter), g.WithClientUserAgent(UserAgent),
		g.WithClientTokenGenerator(c.Generator), g.WithClientTimeout("1m"),
		g.WithClientDialOption(), g.WithClientTLS(c.TLS), g.WithClientID(c.ID),
	}

	if c.Compression {
		opts = append(opts, g.WithClientCompression())
	}

	conn, err := g.NewClient(c.Transport.GRPC.Address, opts...)
	runtime.Must(err)

	return conn
}
