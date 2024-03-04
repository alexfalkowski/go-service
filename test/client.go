package test

import (
	"context"
	"net/http"

	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport"
	tgrpc "github.com/alexfalkowski/go-service/transport/grpc"
	gtracer "github.com/alexfalkowski/go-service/transport/grpc/telemetry/tracer"
	shttp "github.com/alexfalkowski/go-service/transport/http"
	htracer "github.com/alexfalkowski/go-service/transport/http/telemetry/tracer"
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
	tracer, _ := htracer.NewTracer(htracer.Params{Lifecycle: lc, Config: cfg, Version: Version})
	client, _ := shttp.NewClient(
		shttp.WithClientLogger(logger),
		shttp.WithClientRoundTripper(rt), shttp.WithClientBreaker(),
		shttp.WithClientTracer(tracer), shttp.WithClientRetry(&tcfg.HTTP.Retry),
		shttp.WithClientMetrics(meter), shttp.WithClientUserAgent(tcfg.HTTP.UserAgent),
	)

	return client
}

// NewGRPCClient for test.
func NewGRPCClient(
	ctx context.Context, lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *tracer.Config,
	cred credentials.PerRPCCredentials,
	meter metric.Meter,
) *grpc.ClientConn {
	tracer, _ := gtracer.NewTracer(gtracer.Params{Lifecycle: lc, Config: ocfg, Version: Version}) //nolint:contextcheck

	dialOpts := []grpc.DialOption{grpc.WithBlock()}
	if cred != nil {
		dialOpts = append(dialOpts, grpc.WithPerRPCCredentials(cred))
	}

	conn, _ := tgrpc.NewClient(ctx, "127.0.0.1:"+tcfg.GRPC.Port,
		tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
		tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(&tcfg.GRPC.Retry),
		tgrpc.WithClientDialOption(dialOpts...), tgrpc.WithClientMetrics(meter),
		tgrpc.WithClientUserAgent(tcfg.GRPC.UserAgent),
	)

	return conn
}

// NewSecureGRPCClient for test.
func NewSecureGRPCClient(
	ctx context.Context, lc fx.Lifecycle, logger *zap.Logger,
	tcfg *transport.Config, ocfg *tracer.Config,
	meter metric.Meter,
) *grpc.ClientConn {
	tracer, _ := gtracer.NewTracer(gtracer.Params{Lifecycle: lc, Config: ocfg, Version: Version}) //nolint:contextcheck
	sec, _ := tgrpc.WithClientSecure(NewSecureClientConfig())

	conn, _ := tgrpc.NewClient(ctx, "localhost:"+tcfg.GRPC.Port,
		tgrpc.WithClientLogger(logger), tgrpc.WithClientTracer(tracer),
		tgrpc.WithClientBreaker(), tgrpc.WithClientRetry(&tcfg.GRPC.Retry),
		tgrpc.WithClientMetrics(meter), tgrpc.WithClientUserAgent(tcfg.GRPC.UserAgent), sec,
	)

	return conn
}
