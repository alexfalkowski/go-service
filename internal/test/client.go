package test

import (
	"github.com/alexfalkowski/go-service/v2/crypto/tls"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/id"
	"github.com/alexfalkowski/go-service/v2/net"
	"github.com/alexfalkowski/go-service/v2/net/grpc"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
	"github.com/alexfalkowski/go-service/v2/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/telemetry/tracer"
	"github.com/alexfalkowski/go-service/v2/token"
	"github.com/alexfalkowski/go-service/v2/transport"
	g "github.com/alexfalkowski/go-service/v2/transport/grpc"
	gl "github.com/alexfalkowski/go-service/v2/transport/grpc/limiter"
	h "github.com/alexfalkowski/go-service/v2/transport/http"
	hl "github.com/alexfalkowski/go-service/v2/transport/http/limiter"
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
	HTTPLimiter  *hl.Client
	GRPCLimiter  *gl.Client
	Compression  bool
}

// NewTracer for client.
func (c *Client) NewTracer() *tracer.Tracer {
	return NewTracer(c.Lifecycle, c.Tracer)
}

// NewHTTP client for test.
func (c *Client) NewHTTP() *http.Client {
	opts := []h.ClientOption{
		h.WithClientLogger(c.Logger),
		h.WithClientRoundTripper(c.RoundTripper), h.WithClientBreaker(),
		h.WithClientRetry(c.Transport.HTTP.Retry),
		h.WithClientUserAgent(UserAgent),
		h.WithClientTokenGenerator(UserID, c.Generator), h.WithClientTimeout("1m"),
		h.WithClientTLS(c.TLS), h.WithClientID(c.ID),
		h.WithClientLimiter(c.HTTPLimiter),
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
		g.WithClientLimiter(c.GRPCLimiter),
	}

	if c.Compression {
		opts = append(opts, g.WithClientCompression())
	}

	_, target, _ := net.SplitNetworkAddress(c.Transport.GRPC.Address)

	conn, err := g.NewClient(target, opts...)
	runtime.Must(err)

	return conn
}
