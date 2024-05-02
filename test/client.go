package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/client"
	"github.com/alexfalkowski/go-service/runtime"
	"github.com/alexfalkowski/go-service/security"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	g "github.com/alexfalkowski/go-service/transport/grpc"
	h "github.com/alexfalkowski/go-service/transport/http"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client for test.
type Client struct {
	Lifecycle    fx.Lifecycle
	Logger       *zap.Logger
	Tracer       *tracer.Config
	Transport    *transport.Config
	Security     *security.Config
	Credentials  credentials.PerRPCCredentials
	RoundTripper http.RoundTripper
	Meter        metric.Meter
}

// NewHTTP client for test.
func (c *Client) NewHTTP() *http.Client {
	sec, err := h.WithClientSecure(c.Security)
	runtime.Must(err)

	tracer := tracer.NewTracer(c.Lifecycle, Environment, Version, c.Tracer, c.Logger)
	client := h.NewClient(
		h.WithClientLogger(c.Logger),
		h.WithClientRoundTripper(c.RoundTripper), h.WithClientBreaker(),
		h.WithClientTracer(tracer), h.WithClientRetry(c.Transport.HTTP.Retry),
		h.WithClientMetrics(c.Meter), h.WithClientUserAgent(c.Transport.HTTP.UserAgent),
		sec,
	)

	return client
}

func (c *Client) NewGRPC() *grpc.ClientConn {
	tracer := tracer.NewTracer(c.Lifecycle, Environment, Version, c.Tracer, c.Logger)

	dialOpts := []grpc.DialOption{}
	if c.Credentials != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(c.Credentials))
	}

	sec, err := g.WithClientSecure(c.Security)
	runtime.Must(err)

	cl := &client.Config{
		Host:      "localhost:" + c.Transport.GRPC.Port,
		Retry:     c.Transport.GRPC.Retry,
		UserAgent: c.Transport.GRPC.UserAgent,
	}

	conn, err := g.NewClient(cl.Host,
		g.WithClientUnaryInterceptors(), g.WithClientStreamInterceptors(),
		g.WithClientLogger(c.Logger), g.WithClientTracer(tracer),
		g.WithClientBreaker(), g.WithClientRetry(cl.Retry),
		g.WithClientDialOption(dialOpts...), g.WithClientMetrics(c.Meter),
		g.WithClientUserAgent(cl.UserAgent), sec,
	)
	runtime.Must(err)

	return conn
}
