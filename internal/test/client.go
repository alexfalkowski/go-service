package test

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/limiter"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	g "github.com/alexfalkowski/go-service/v2/transport/grpc"
	h "github.com/alexfalkowski/go-service/v2/transport/http"
)

// Client for test.
type Client struct {
	Lifecycle    di.Lifecycle
	Logger       *logger.Logger
	Tracer       *tracer.Config
	Transport    *transport.Config
	TLS          *tls.Config
	ID           id.Generator
	RoundTripper http.RoundTripper
	Meter        *metrics.Meter
	Generator    token.Generator
	Limiter      *limiter.Limiter
	Compression  bool
}

// NewTracer for client.
func (c *Client) NewTracer() *tracer.Tracer {
	return NewTracer(c.Lifecycle, c.Tracer)
}

// NewHTTP client for test.
func (c *Client) NewHTTP() *http.Client {
	tracer := c.NewTracer()
	opts := []h.ClientOption{
		h.WithClientLogger(c.Logger),
		h.WithClientRoundTripper(c.RoundTripper), h.WithClientBreaker(),
		h.WithClientTracer(tracer), h.WithClientRetry(c.Transport.HTTP.Retry),
		h.WithClientMetrics(c.Meter), h.WithClientUserAgent(UserAgent),
		h.WithClientTokenGenerator(UserID, c.Generator), h.WithClientTimeout("1m"),
		h.WithClientTLS(c.TLS), h.WithClientID(c.ID), h.WithClientLimiter(c.Limiter),
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
		g.WithClientTokenGenerator(UserID, c.Generator), g.WithClientTimeout("1m"),
		g.WithClientDialOption(), g.WithClientTLS(c.TLS), g.WithClientID(c.ID),
		g.WithClientLimiter(c.Limiter),
	}

	if c.Compression {
		opts = append(opts, g.WithClientCompression())
	}

	conn, err := g.NewClient(c.Transport.GRPC.Address, opts...)
	runtime.Must(err)

	return conn
}
