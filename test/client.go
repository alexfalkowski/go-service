package test

import (
	"net/http"

	"github.com/alexfalkowski/go-service/client"
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

// NewHTTPClient for test.
func NewHTTPClient(lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config, meter metric.Meter) *http.Client {
	return NewHTTPClientWithRoundTripper(lc, logger, cfg, tcfg, nil, meter)
}

// NewHTTPClientWithRoundTripper for test.
func NewHTTPClientWithRoundTripper(lc fx.Lifecycle, logger *zap.Logger, cfg *tracer.Config, tcfg *transport.Config, rt http.RoundTripper, meter metric.Meter) *http.Client {
	tracer, err := tracer.NewTracer(lc, Environment, Version, cfg)
	if err != nil {
		panic(err)
	}

	client := h.NewClient(
		h.WithClientLogger(logger),
		h.WithClientRoundTripper(rt), h.WithClientBreaker(),
		h.WithClientTracer(tracer), h.WithClientRetry(tcfg.HTTP.Retry),
		h.WithClientMetrics(meter), h.WithClientUserAgent(tcfg.HTTP.UserAgent),
	)

	return client
}

// NewGRPCClient for test.
func NewGRPCClient(
	lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *tracer.Config,
	cred credentials.PerRPCCredentials,
	meter metric.Meter,
) *grpc.ClientConn {
	tracer, err := tracer.NewTracer(lc, Environment, Version, ocfg)
	if err != nil {
		panic(err)
	}

	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	if cred != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(cred))
	}

	cl := &client.Config{Host: "127.0.0.1:" + tcfg.GRPC.Port, Retry: tcfg.GRPC.Retry, UserAgent: tcfg.GRPC.UserAgent}

	conn, err := g.NewClient(cl.Host,
		g.WithClientLogger(logger), g.WithClientTracer(tracer),
		g.WithClientBreaker(), g.WithClientRetry(cl.Retry),
		g.WithClientDialOption(dialOpts...), g.WithClientMetrics(meter),
		g.WithClientUserAgent(cl.UserAgent),
	)
	if err != nil {
		panic(err)
	}

	return conn
}

// NewSecureGRPCClient for test.
func NewSecureGRPCClient(
	lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *tracer.Config,
	meter metric.Meter,
) *grpc.ClientConn {
	tracer, err := tracer.NewTracer(lc, Environment, Version, ocfg)
	if err != nil {
		panic(err)
	}

	sec, err := g.WithClientSecure(NewSecureClientConfig())
	if err != nil {
		panic(err)
	}

	conn, err := g.NewClient("localhost:"+tcfg.GRPC.Port,
		g.WithClientLogger(logger), g.WithClientTracer(tracer),
		g.WithClientBreaker(), g.WithClientRetry(tcfg.GRPC.Retry),
		g.WithClientMetrics(meter), g.WithClientUserAgent(tcfg.GRPC.UserAgent), sec,
	)
	if err != nil {
		panic(err)
	}

	return conn
}
